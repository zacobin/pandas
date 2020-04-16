// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"crypto/aes"
	"encoding/hex"
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
	rediscons "github.com/cloustone/pandas/bootstrap/redis/consumer"
	redisprod "github.com/cloustone/pandas/bootstrap/redis/producer"
	"github.com/cloustone/pandas/pkg/logger"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/cloustone/pandas/bootstrap"
	api "github.com/cloustone/pandas/bootstrap/api"
	"github.com/cloustone/pandas/bootstrap/postgres"
	"github.com/cloustone/pandas/mainflux"
	mflog "github.com/cloustone/pandas/pkg/logger"
	mfsdk "github.com/cloustone/pandas/sdk/go"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	r "github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defLogLevel       = "error"
	defDBHost         = "localhost"
	defDBPort         = "5432"
	defDBUser         = "mainflux"
	defDBPass         = "mainflux"
	defDBName         = "bootstrap"
	defDBSSLMode      = "disable"
	defDBSSLCert      = ""
	defDBSSLKey       = ""
	defDBSSLRootCert  = ""
	defEncryptKey     = "12345678910111213141516171819202"
	defClientTLS      = "false"
	defCACerts        = ""
	defPort           = "8180"
	defServerCert     = ""
	defServerKey      = ""
	defBaseURL        = "http://localhost"
	defThingsPrefix   = ""
	defThingsESURL    = "localhost:6379"
	defThingsESPass   = ""
	defThingsESDB     = "0"
	defESURL          = "localhost:6379"
	defESPass         = ""
	defESDB           = "0"
	defESConsumerName = "bootstrap"
	defJaegerURL      = ""
	defAuthURL        = "localhost:8181"
	defAuthTimeout    = "1" // in seconds

	envLogLevel       = "PD_BOOTSTRAP_LOG_LEVEL"
	envDBHost         = "PD_BOOTSTRAP_DB_HOST"
	envDBPort         = "PD_BOOTSTRAP_DB_PORT"
	envDBUser         = "PD_BOOTSTRAP_DB_USER"
	envDBPass         = "PD_BOOTSTRAP_DB_PASS"
	envDBName         = "PD_BOOTSTRAP_DB"
	envDBSSLMode      = "PD_BOOTSTRAP_DB_SSL_MODE"
	envDBSSLCert      = "PD_BOOTSTRAP_DB_SSL_CERT"
	envDBSSLKey       = "PD_BOOTSTRAP_DB_SSL_KEY"
	envDBSSLRootCert  = "PD_BOOTSTRAP_DB_SSL_ROOT_CERT"
	envEncryptKey     = "PD_BOOTSTRAP_ENCRYPT_KEY"
	envClientTLS      = "PD_BOOTSTRAP_CLIENT_TLS"
	envCACerts        = "PD_BOOTSTRAP_CA_CERTS"
	envPort           = "PD_BOOTSTRAP_PORT"
	envServerCert     = "PD_BOOTSTRAP_SERVER_CERT"
	envServerKey      = "PD_BOOTSTRAP_SERVER_KEY"
	envBaseURL        = "PD_SDK_BASE_URL"
	envThingsPrefix   = "PD_SDK_THINGS_PREFIX"
	envThingsESURL    = "PD_THINGS_ES_URL"
	envThingsESPass   = "PD_THINGS_ES_PASS"
	envThingsESDB     = "PD_THINGS_ES_DB"
	envESURL          = "PD_BOOTSTRAP_ES_URL"
	envESPass         = "PD_BOOTSTRAP_ES_PASS"
	envESDB           = "PD_BOOTSTRAP_ES_DB"
	envESConsumerName = "PD_BOOTSTRAP_EVENT_CONSUMER"
	envJaegerURL      = "PD_JAEGER_URL"
	envAuthURL        = "PD_AUTH_URL"
	envAuthTimeout    = "PD_AUTH_TIMEOUT"
)

type config struct {
	logLevel       string
	dbConfig       postgres.Config
	clientTLS      bool
	encKey         []byte
	caCerts        string
	httpPort       string
	serverCert     string
	serverKey      string
	baseURL        string
	thingsPrefix   string
	esThingsURL    string
	esThingsPass   string
	esThingsDB     string
	esURL          string
	esPass         string
	esDB           string
	esConsumerName string
	jaegerURL      string
	authURL        string
	authTimeout    time.Duration
}

func main() {
	cfg := loadConfig()

	logger, err := mflog.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	thingsESConn := connectToRedis(cfg.esThingsURL, cfg.esThingsPass, cfg.esThingsDB, logger)
	defer thingsESConn.Close()

	esClient := connectToRedis(cfg.esURL, cfg.esPass, cfg.esDB, logger)
	defer esClient.Close()

	authTracer, authCloser := initJaeger("auth", cfg.jaegerURL, logger)
	defer authCloser.Close()

	authConn := connectToAuth(cfg, logger)
	defer authConn.Close()

	auth := authapi.NewClient(authTracer, authConn, cfg.authTimeout)

	svc := newService(auth, db, logger, esClient, cfg)
	errs := make(chan error, 2)

	go startHTTPServer(svc, cfg, logger, errs)
	go subscribeToThingsES(svc, thingsESConn, cfg.esConsumerName, logger)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Bootstrap service terminated: %s", err))
}

func loadConfig() config {
	tls, err := strconv.ParseBool(pandas.Env(envClientTLS, defClientTLS))
	if err != nil {
		tls = false
	}
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

	timeout, err := strconv.ParseInt(pandas.Env(envAuthTimeout, defAuthTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envAuthTimeout, err.Error())
	}
	encKey, err := hex.DecodeString(pandas.Env(envEncryptKey, defEncryptKey))
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envEncryptKey, err.Error())
	}
	if err := os.Unsetenv(envEncryptKey); err != nil {
		log.Fatalf("Unable to unset %s value: %s", envEncryptKey, err.Error())
	}
	if _, err := aes.NewCipher(encKey); err != nil {
		log.Fatalf("Invalid %s value: %s", envEncryptKey, err.Error())
	}

	return config{
		logLevel:       pandas.Env(envLogLevel, defLogLevel),
		dbConfig:       dbConfig,
		clientTLS:      tls,
		encKey:         encKey,
		caCerts:        pandas.Env(envCACerts, defCACerts),
		httpPort:       pandas.Env(envPort, defPort),
		serverCert:     pandas.Env(envServerCert, defServerCert),
		serverKey:      pandas.Env(envServerKey, defServerKey),
		baseURL:        pandas.Env(envBaseURL, defBaseURL),
		thingsPrefix:   pandas.Env(envThingsPrefix, defThingsPrefix),
		esThingsURL:    pandas.Env(envThingsESURL, defThingsESURL),
		esThingsPass:   pandas.Env(envThingsESPass, defThingsESPass),
		esThingsDB:     pandas.Env(envThingsESDB, defThingsESDB),
		esURL:          pandas.Env(envESURL, defESURL),
		esPass:         pandas.Env(envESPass, defESPass),
		esDB:           pandas.Env(envESDB, defESDB),
		esConsumerName: pandas.Env(envESConsumerName, defESConsumerName),
		jaegerURL:      pandas.Env(envJaegerURL, defJaegerURL),
		authURL:        pandas.Env(envAuthURL, defAuthURL),
		authTimeout:    time.Duration(timeout) * time.Second,
	}
}

func connectToDB(cfg postgres.Config, logger mflog.Logger) *sqlx.DB {
	db, err := postgres.Connect(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to postgres: %s", err))
		os.Exit(1)
	}
	return db
}

func connectToRedis(redisURL, redisPass, redisDB string, logger mflog.Logger) *r.Client {
	db, err := strconv.Atoi(redisDB)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to redis: %s", err))
		os.Exit(1)
	}

	return r.NewClient(&r.Options{
		Addr:     redisURL,
		Password: redisPass,
		DB:       db,
	})
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

func newService(auth mainflux.AuthNServiceClient, db *sqlx.DB, logger mflog.Logger, esClient *r.Client, cfg config) bootstrap.Service {
	thingsRepo := postgres.NewConfigRepository(db, logger)

	config := mfsdk.Config{
		BaseURL:      cfg.baseURL,
		ThingsPrefix: cfg.thingsPrefix,
	}

	sdk := mfsdk.NewSDK(config)

	svc := bootstrap.New(auth, thingsRepo, sdk, cfg.encKey)
	svc = redisprod.NewEventStoreMiddleware(svc, esClient)
	svc = api.NewLoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "bootstrap",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "bootstrap",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
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

	conn, err := grpc.Dial(cfg.authURL, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to users service: %s", err))
		os.Exit(1)
	}

	return conn
}

func startHTTPServer(svc bootstrap.Service, cfg config, logger mflog.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.httpPort)
	if cfg.serverCert != "" || cfg.serverKey != "" {
		logger.Info(fmt.Sprintf("Bootstrap service started using https on port %s with cert %s key %s",
			cfg.httpPort, cfg.serverCert, cfg.serverKey))
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, api.MakeHandler(svc, bootstrap.NewConfigReader(cfg.encKey)))
		return
	}
	logger.Info(fmt.Sprintf("Bootstrap service started using http on port %s", cfg.httpPort))
	errs <- http.ListenAndServe(p, api.MakeHandler(svc, bootstrap.NewConfigReader(cfg.encKey)))
}

func subscribeToThingsES(svc bootstrap.Service, client *r.Client, consumer string, logger mflog.Logger) {
	eventStore := rediscons.NewEventStore(svc, client, consumer, logger)
	logger.Info("Subscribed to Redis Event Store")
	if err := eventStore.Subscribe("mainflux.things"); err != nil {
		logger.Warn(fmt.Sprintf("Bootstrap service failed to subscribe to event sourcing: %s", err))
	}
}
