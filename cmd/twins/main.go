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
	authapi "github.com/cloustone/pandas/authn/api/grpc"
	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/logger"
	localusers "github.com/cloustone/pandas/things/users"
	"github.com/cloustone/pandas/twins"
	"github.com/cloustone/pandas/twins/api"
	twapi "github.com/cloustone/pandas/twins/api/http"
	twmongodb "github.com/cloustone/pandas/twins/mongodb"
	natspub "github.com/cloustone/pandas/twins/nats/publisher"
	natssub "github.com/cloustone/pandas/twins/nats/subscriber"
	"github.com/cloustone/pandas/twins/uuid"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	nats "github.com/nats-io/nats.go"
	opentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"go.mongodb.org/mongo-driver/mongo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defLogLevel        = "info"
	defHTTPPort        = "9021"
	defJaegerURL       = ""
	defServerCert      = ""
	defServerKey       = ""
	defDBName          = "mainflux"
	defDBHost          = "localhost"
	defDBPort          = "27017"
	defSingleUserEmail = ""
	defSingleUserToken = ""
	defClientTLS       = "false"
	defCACerts         = ""
	defThingID         = ""
	defThingKey        = ""
	defChannelID       = ""
	defNatsURL         = nats.DefaultURL

	defAuthnTimeout = "1" // in seconds
	defAuthnURL     = "localhost:8181"

	envLogLevel        = "PD_TWINS_LOG_LEVEL"
	envHTTPPort        = "PD_TWINS_HTTP_PORT"
	envJaegerURL       = "PD_JAEGER_URL"
	envServerCert      = "PD_TWINS_SERVER_CERT"
	envServerKey       = "PD_TWINS_SERVER_KEY"
	envDBName          = "PD_TWINS_DB_NAME"
	envDBHost          = "PD_TWINS_DB_HOST"
	envDBPort          = "PD_TWINS_DB_PORT"
	envSingleUserEmail = "PD_TWINS_SINGLE_USER_EMAIL"
	envSingleUserToken = "PD_TWINS_SINGLE_USER_TOKEN"
	envClientTLS       = "PD_TWINS_CLIENT_TLS"
	envCACerts         = "PD_TWINS_CA_CERTS"
	envThingID         = "PD_TWINS_THING_ID"
	envThingKey        = "PD_TWINS_THING_KEY"
	envChannelID       = "PD_TWINS_CHANNEL_ID"
	envNatsURL         = "PD_NATS_URL"

	envAuthnTimeout = "PD_AUTHN_TIMEOUT"
	envAuthnURL     = "PD_AUTHN_URL"
)

type config struct {
	logLevel        string
	httpPort        string
	jaegerURL       string
	serverCert      string
	serverKey       string
	dbCfg           twmongodb.Config
	singleUserEmail string
	singleUserToken string
	clientTLS       bool
	caCerts         string
	thingID         string
	thingKey        string
	channelID       string
	NatsURL         string

	authnTimeout time.Duration
	authnURL     string
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db, err := twmongodb.Connect(cfg.dbCfg, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	authTracer, authCloser := initJaeger("auth", cfg.jaegerURL, logger)
	defer authCloser.Close()

	auth, _ := createAuthClient(cfg, authTracer, logger)

	dbTracer, dbCloser := initJaeger("twins_db", cfg.jaegerURL, logger)
	defer dbCloser.Close()

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to NATS: %s", err))
		os.Exit(1)
	}
	defer nc.Close()

	ncTracer, ncCloser := initJaeger("twins_nats", cfg.jaegerURL, logger)
	defer ncCloser.Close()

	tracer, closer := initJaeger("twins", cfg.jaegerURL, logger)
	defer closer.Close()

	svc := newService(nc, ncTracer, cfg.channelID, auth, dbTracer, db, logger)

	errs := make(chan error, 2)

	go startHTTPServer(twapi.MakeHandler(tracer, svc), cfg.httpPort, cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Twins service terminated: %s", err))
}

func loadConfig() config {
	tls, err := strconv.ParseBool(pandas.Env(envClientTLS, defClientTLS))
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envClientTLS)
	}

	timeout, err := strconv.ParseInt(pandas.Env(envAuthnTimeout, defAuthnTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envAuthnTimeout, err.Error())
	}

	dbCfg := twmongodb.Config{
		Name: pandas.Env(envDBName, defDBName),
		Host: pandas.Env(envDBHost, defDBHost),
		Port: pandas.Env(envDBPort, defDBPort),
	}

	return config{
		logLevel:        pandas.Env(envLogLevel, defLogLevel),
		httpPort:        pandas.Env(envHTTPPort, defHTTPPort),
		serverCert:      pandas.Env(envServerCert, defServerCert),
		serverKey:       pandas.Env(envServerKey, defServerKey),
		jaegerURL:       pandas.Env(envJaegerURL, defJaegerURL),
		dbCfg:           dbCfg,
		singleUserEmail: pandas.Env(envSingleUserEmail, defSingleUserEmail),
		singleUserToken: pandas.Env(envSingleUserToken, defSingleUserToken),
		clientTLS:       tls,
		caCerts:         pandas.Env(envCACerts, defCACerts),
		thingID:         pandas.Env(envThingID, defThingID),
		channelID:       pandas.Env(envChannelID, defChannelID),
		thingKey:        pandas.Env(envThingKey, defThingKey),
		NatsURL:         pandas.Env(envNatsURL, defNatsURL),
		authnURL:        pandas.Env(envAuthnURL, defAuthnURL),
		authnTimeout:    time.Duration(timeout) * time.Second,
	}
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

func createAuthClient(cfg config, tracer opentracing.Tracer, logger logger.Logger) (mainflux.AuthNServiceClient, func() error) {
	if cfg.singleUserEmail != "" && cfg.singleUserToken != "" {
		return localusers.NewSingleUserService(cfg.singleUserEmail, cfg.singleUserToken), nil
	}

	conn := connectToAuth(cfg, logger)
	return authapi.NewClient(tracer, conn, cfg.authnTimeout), conn.Close
}

func connectToAuth(cfg config, logger logger.Logger) *grpc.ClientConn {
	var opts []grpc.DialOption
	if cfg.clientTLS {
		if cfg.caCerts != "" {
			tpc, err := credentials.NewClientTLSFromFile(cfg.caCerts, "")
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to create tls credentials: %s", err))
				os.Exit(1)
			}
			opts = append(opts, grpc.WithTransportCredentials(tpc))
		}
	} else {
		opts = append(opts, grpc.WithInsecure())
		logger.Info("gRPC communication is not encrypted")
	}

	conn, err := grpc.Dial(cfg.authnURL, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to auth service: %s", err))
		os.Exit(1)
	}

	return conn
}

func newService(nc *nats.Conn, ncTracer opentracing.Tracer, chanID string, users mainflux.AuthNServiceClient, dbTracer opentracing.Tracer, db *mongo.Database, logger logger.Logger) twins.Service {
	twinRepo := twmongodb.NewTwinRepository(db)

	stateRepo := twmongodb.NewStateRepository(db)
	idp := uuid.New()

	np := natspub.NewPublisher(nc, chanID, logger)

	svc := twins.New(users, twinRepo, stateRepo, idp, np)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "twins",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "twins",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	natssub.NewSubscriber(nc, chanID, svc, logger)

	return svc
}

func startHTTPServer(handler http.Handler, port string, cfg config, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if cfg.serverCert != "" || cfg.serverKey != "" {
		logger.Info(fmt.Sprintf("Twins service started using https on port %s with cert %s key %s",
			port, cfg.serverCert, cfg.serverKey))
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, handler)
		return
	}
	logger.Info(fmt.Sprintf("Twins service started using http on port %s", cfg.httpPort))
	errs <- http.ListenAndServe(p, handler)
}
