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
	"github.com/cloustone/pandas/pms"
	"github.com/cloustone/pandas/pms/tracing"

	"github.com/jmoiron/sqlx"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/credentials"

	authapi "github.com/cloustone/pandas/authn/api/grpc"
	"github.com/cloustone/pandas/pkg/logger"
	"github.com/cloustone/pandas/pms/api"
	httpapi "github.com/cloustone/pandas/pms/api/http"
	"github.com/cloustone/pandas/pms/postgres"
	rediscache "github.com/cloustone/pandas/pms/redis"
	"github.com/cloustone/pandas/pms/uuid"
	localusers "github.com/cloustone/pandas/things/users"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-redis/redis"
	nats "github.com/nats-io/nats.go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"

	natspub "github.com/cloustone/pandas/pms/nats/publisher"
	natssub "github.com/cloustone/pandas/pms/nats/subscriber"
)

const (
	defLogLevel        = "error"
	defDBHost          = "localhost"
	defDBPort          = "5432"
	defDBUser          = "mainflux"
	defDBPass          = "mainflux"
	defDBName          = "things"
	defDBSSLMode       = "disable"
	defDBSSLCert       = ""
	defDBSSLKey        = ""
	defDBSSLRootCert   = ""
	defClientTLS       = "false"
	defCACerts         = ""
	defCacheURL        = "localhost:6379"
	defCachePass       = ""
	defCacheDB         = "0"
	defESURL           = "localhost:6379"
	defESPass          = ""
	defESDB            = "0"
	defHTTPPort        = "8180"
	defAuthHTTPPort    = "8989"
	defAuthGRPCPort    = "8181"
	defServerCert      = ""
	defServerKey       = ""
	defSingleUserEmail = ""
	defSingleUserToken = ""
	defJaegerURL       = ""
	defAuthURL         = "localhost:8181"
	defAuthTimeout     = "1" // in seconds
	defNatsURL         = nats.DefaultURL
	defChannelID       = ""

	envLogLevel        = "PD_V2MS_LOG_LEVEL"
	envDBHost          = "PD_V2MS_DB_HOST"
	envDBPort          = "PD_V2MS_DB_PORT"
	envDBUser          = "PD_V2MS_DB_USER"
	envDBPass          = "PD_V2MS_DB_PASS"
	envDBName          = "PD_V2MS_DB"
	envDBSSLMode       = "PD_V2MS_DB_SSL_MODE"
	envDBSSLCert       = "PD_V2MS_DB_SSL_CERT"
	envDBSSLKey        = "PD_V2MS_DB_SSL_KEY"
	envDBSSLRootCert   = "PD_V2MS_DB_SSL_ROOT_CERT"
	envClientTLS       = "PD_V2MS_CLIENT_TLS"
	envCACerts         = "PD_V2MS_CA_CERTS"
	envCacheURL        = "PD_V2MS_CACHE_URL"
	envCachePass       = "PD_V2MS_CACHE_PASS"
	envCacheDB         = "PD_V2MS_CACHE_DB"
	envESURL           = "PD_V2MS_ES_URL"
	envESPass          = "PD_V2MS_ES_PASS"
	envESDB            = "PD_V2MS_ES_DB"
	envHTTPPort        = "PD_V2MS_HTTP_PORT"
	envAuthHTTPPort    = "PD_V2MS_AUTH_HTTP_PORT"
	envAuthGRPCPort    = "PD_V2MS_AUTH_GRPC_PORT"
	envServerCert      = "PD_V2MS_SERVER_CERT"
	envServerKey       = "PD_V2MS_SERVER_KEY"
	envSingleUserEmail = "PD_V2MS_SINGLE_USER_EMAIL"
	envSingleUserToken = "PD_V2MS_SINGLE_USER_TOKEN"
	envJaegerURL       = "PD_JAEGER_URL"
	envAuthURL         = "PD_AUTH_URL"
	envAuthTimeout     = "PD_AUTH_TIMEOUT"
	envNatsURL         = "PD_NATS_URL"
	envChannelID       = "PD_V2MS_CHANNEL_ID"
)

type config struct {
	logLevel        string
	dbConfig        postgres.Config
	clientTLS       bool
	caCerts         string
	cacheURL        string
	cachePass       string
	cacheDB         string
	esURL           string
	esPass          string
	esDB            string
	httpPort        string
	authHTTPPort    string
	authGRPCPort    string
	serverCert      string
	serverKey       string
	singleUserEmail string
	singleUserToken string
	jaegerURL       string
	authURL         string
	authTimeout     time.Duration
	NatsURL         string
	channelID       string
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	thingsTracer, thingsCloser := initJaeger("things", cfg.jaegerURL, logger)
	defer thingsCloser.Close()

	cacheClient := connectToRedis(cfg.cacheURL, cfg.cachePass, cfg.cacheDB, logger)

	esClient := connectToRedis(cfg.esURL, cfg.esPass, cfg.esDB, logger)

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	authTracer, authCloser := initJaeger("auth", cfg.jaegerURL, logger)
	defer authCloser.Close()

	auth, close := createAuthClient(cfg, authTracer, logger)
	if close != nil {
		defer close()
	}

	dbTracer, dbCloser := initJaeger("pms_db", cfg.jaegerURL, logger)
	defer dbCloser.Close()

	cacheTracer, cacheCloser := initJaeger("pms_cache", cfg.jaegerURL, logger)
	defer cacheCloser.Close()

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to NATS: %s", err))
		os.Exit(1)
	}
	defer nc.Close()

	ncTracer, ncCloser := initJaeger("pms_nats", cfg.jaegerURL, logger)
	defer ncCloser.Close()

	svc := newService(nc, ncTracer, cfg.channelID, auth, dbTracer, cacheTracer, db, cacheClient, esClient, logger)
	errs := make(chan error, 2)

	go startHTTPServer(httpapi.MakeHandler(thingsTracer, svc), cfg.httpPort, cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("PMS service terminated: %s", err))
}

func loadConfig() config {
	tls, err := strconv.ParseBool(pandas.Env(envClientTLS, defClientTLS))
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envClientTLS)
	}

	timeout, err := strconv.ParseInt(pandas.Env(envAuthTimeout, defAuthTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envAuthTimeout, err.Error())
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

	return config{
		logLevel:        pandas.Env(envLogLevel, defLogLevel),
		dbConfig:        dbConfig,
		clientTLS:       tls,
		caCerts:         pandas.Env(envCACerts, defCACerts),
		cacheURL:        pandas.Env(envCacheURL, defCacheURL),
		cachePass:       pandas.Env(envCachePass, defCachePass),
		cacheDB:         pandas.Env(envCacheDB, defCacheDB),
		esURL:           pandas.Env(envESURL, defESURL),
		esPass:          pandas.Env(envESPass, defESPass),
		esDB:            pandas.Env(envESDB, defESDB),
		httpPort:        pandas.Env(envHTTPPort, defHTTPPort),
		authHTTPPort:    pandas.Env(envAuthHTTPPort, defAuthHTTPPort),
		authGRPCPort:    pandas.Env(envAuthGRPCPort, defAuthGRPCPort),
		serverCert:      pandas.Env(envServerCert, defServerCert),
		serverKey:       pandas.Env(envServerKey, defServerKey),
		singleUserEmail: pandas.Env(envSingleUserEmail, defSingleUserEmail),
		singleUserToken: pandas.Env(envSingleUserToken, defSingleUserToken),
		jaegerURL:       pandas.Env(envJaegerURL, defJaegerURL),
		authURL:         pandas.Env(envAuthURL, defAuthURL),
		NatsURL:         pandas.Env(envNatsURL, defNatsURL),
		channelID:       pandas.Env(envChannelID, defChannelID),
		authTimeout:     time.Duration(timeout) * time.Second,
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

func connectToRedis(cacheURL, cachePass string, cacheDB string, logger logger.Logger) *redis.Client {
	db, err := strconv.Atoi(cacheDB)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to cache: %s", err))
		os.Exit(1)
	}

	return redis.NewClient(&redis.Options{
		Addr:     cacheURL,
		Password: cachePass,
		DB:       db,
	})
}

func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sqlx.DB {
	db, err := postgres.Connect(dbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to postgres: %s", err))
		os.Exit(1)
	}
	return db
}

func createAuthClient(cfg config, tracer opentracing.Tracer, logger logger.Logger) (mainflux.AuthNServiceClient, func() error) {
	if cfg.singleUserEmail != "" && cfg.singleUserToken != "" {
		return localusers.NewSingleUserService(cfg.singleUserEmail, cfg.singleUserToken), nil
	}

	conn := connectToAuth(cfg, logger)
	return authapi.NewClient(tracer, conn, cfg.authTimeout), conn.Close
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

func newService(nc *nats.Conn, ncTracer opentracing.Tracer, chanID string, auth mainflux.AuthNServiceClient, dbTracer opentracing.Tracer, cacheTracer opentracing.Tracer, db *sqlx.DB, cacheClient *redis.Client, esClient *redis.Client, logger logger.Logger) pms.Service {
	database := postgres.NewDatabase(db)

	projectsRepo := postgres.NewProjectRepository(database)
	projectsRepo = tracing.ProjectRepositoryMiddleware(dbTracer, projectsRepo)

	projectCache := rediscache.NewProjectCache(cacheClient)
	projectCache = tracing.ProjectCacheMiddleware(cacheTracer, projectCache)

	idp := uuid.New()

	np := natspub.NewPublisher(nc, chanID, logger)
	svc := pms.New(auth, projectsRepo, projectCache, idp, np)
	//svc = rediscache.NewEventStoreMiddleware(svc, esClient)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "things",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "things",
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
		logger.Info(fmt.Sprintf("pms service started using https on port %s with cert %s key %s",
			port, cfg.serverCert, cfg.serverKey))
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, handler)
		return
	}
	logger.Info(fmt.Sprintf("pms service started using http on port %s", cfg.httpPort))
	errs <- http.ListenAndServe(p, handler)
}
