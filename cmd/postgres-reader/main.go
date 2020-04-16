// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/mainflux/readers"
	"github.com/cloustone/pandas/mainflux/readers/api"
	"github.com/cloustone/pandas/mainflux/readers/postgres"
	"github.com/cloustone/pandas/pkg/logger"
	thingsapi "github.com/cloustone/pandas/things/api/auth/grpc"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	opentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	svcName = "postgres-writer"
	sep     = ","

	defThingsURL     = "localhost:8183"
	defLogLevel      = "debug"
	defPort          = "9204"
	defClientTLS     = "false"
	defCACerts       = ""
	defDBHost        = "localhost"
	defDBPort        = "5432"
	defDBUser        = "mainflux"
	defDBPass        = "mainflux"
	defDBName        = "messages"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defJaegerURL     = ""
	defThingsTimeout = "1" // in seconds

	envThingsURL     = "PD_THINGS_URL"
	envLogLevel      = "PD_POSTGRES_READER_LOG_LEVEL"
	envPort          = "PD_POSTGRES_READER_PORT"
	envClientTLS     = "PD_POSTGRES_READER_CLIENT_TLS"
	envCACerts       = "PD_POSTGRES_READER_CA_CERTS"
	envDBHost        = "PD_POSTGRES_READER_DB_HOST"
	envDBPort        = "PD_POSTGRES_READER_DB_PORT"
	envDBUser        = "PD_POSTGRES_READER_DB_USER"
	envDBPass        = "PD_POSTGRES_READER_DB_PASS"
	envDBName        = "PD_POSTGRES_READER_DB_NAME"
	envDBSSLMode     = "PD_POSTGRES_READER_DB_SSL_MODE"
	envDBSSLCert     = "PD_POSTGRES_READER_DB_SSL_CERT"
	envDBSSLKey      = "PD_POSTGRES_READER_DB_SSL_KEY"
	envDBSSLRootCert = "PD_POSTGRES_READER_DB_SSL_ROOT_CERT"
	envJaegerURL     = "PD_JAEGER_URL"
	envThingsTimeout = "PD_POSTGRES_READER_THINGS_TIMEOUT"
)

type config struct {
	thingsURL     string
	logLevel      string
	port          string
	clientTLS     bool
	caCerts       string
	dbConfig      postgres.Config
	jaegerURL     string
	thingsTimeout time.Duration
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	conn := connectToThings(cfg, logger)
	defer conn.Close()

	thingsTracer, thingsCloser := initJaeger("things", cfg.jaegerURL, logger)
	defer thingsCloser.Close()

	tc := thingsapi.NewClient(conn, thingsTracer, cfg.thingsTimeout)

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	repo := newService(db, logger)

	errs := make(chan error, 2)

	go startHTTPServer(repo, tc, cfg.port, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Postgres writer service terminated: %s", err))
}

func loadConfig() config {
	dbConfig := postgres.Config{
		Host:        pandas.Env(envDBHost, defDBHost),
		Port:        pandas.Env(envDBPort, defDBPort),
		User:        pandas.Env(envDBUser, defDBUser),
		Pass:        pandas.Env(envDBPass, defDBPass),
		Name:        pandas.Env(envDBName, defDBName),
		SSLMode:     pandas.Env(envDBSSLMode, defDBSSLMode),
		SSLCert:     pandas.Env(envDBSSLCert, defDBSSLCert),
		SSLKey:      pandas.Env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: pandas.Env(envDBSSLRootCert, defDBSSLRootCert),
	}

	timeout, err := strconv.ParseInt(pandas.Env(envThingsTimeout, defThingsTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envThingsTimeout, err.Error())
	}

	return config{
		thingsURL:     pandas.Env(envThingsURL, defThingsURL),
		logLevel:      pandas.Env(envLogLevel, defLogLevel),
		port:          pandas.Env(envPort, defPort),
		dbConfig:      dbConfig,
		jaegerURL:     pandas.Env(envJaegerURL, defJaegerURL),
		thingsTimeout: time.Duration(timeout) * time.Second,
	}
}

func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sqlx.DB {
	db, err := postgres.Connect(dbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to Postgres: %s", err))
		os.Exit(1)
	}
	return db
}

func initJaeger(svcName, url string, logger logger.Logger) (opentracing.Tracer, io.Closer) {
	if url == "" {
		return opentracing.NoopTracer{}, ioutil.NopCloser(nil)
	}

	tracer, closer, err := jconfig.Configuration{
		ServiceName: svcName,
		Sampler: &jconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: url,
			LogSpans:           true,
		},
	}.NewTracer()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to init Jaeger client: %s", err))
		os.Exit(1)
	}

	return tracer, closer
}

func connectToThings(cfg config, logger logger.Logger) *grpc.ClientConn {
	var opts []grpc.DialOption
	if cfg.clientTLS {
		if cfg.caCerts != "" {
			tpc, err := credentials.NewClientTLSFromFile(cfg.caCerts, "")
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to load certs: %s", err))
				os.Exit(1)
			}
			opts = append(opts, grpc.WithTransportCredentials(tpc))
		}
	} else {
		logger.Info("gRPC communication is not encrypted")
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(cfg.thingsURL, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to things service: %s", err))
		os.Exit(1)
	}
	return conn
}

func newService(db *sqlx.DB, logger logger.Logger) readers.MessageRepository {
	svc := postgres.New(db)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "postgres",
			Subsystem: "message_writer",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "postgres",
			Subsystem: "message_writer",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}

func startHTTPServer(repo readers.MessageRepository, tc mainflux.ThingsServiceClient, port string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Postgres reader service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHandler(repo, tc, svcName))
}
