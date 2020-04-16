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
	"github.com/cloustone/pandas/pkg/email"
	"github.com/cloustone/pandas/rulechain"
	"github.com/cloustone/pandas/rulechain/tracing"
	"github.com/go-redis/redis"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	authapi "github.com/cloustone/pandas/authn/api/grpc"
	"github.com/cloustone/pandas/pkg/logger"
	"github.com/cloustone/pandas/rulechain/api"
	rulechainapi "github.com/cloustone/pandas/rulechain/api/http"
	natssub "github.com/cloustone/pandas/rulechain/nats/subscriber"
	"github.com/cloustone/pandas/rulechain/postgres"
	rediscache "github.com/cloustone/pandas/rulechain/redis"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	opentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
)

const (
	defLogLevel      = "error"
	defDBHost        = "localhost"
	defDBPort        = "5432"
	defDBUser        = "mainflux"
	defDBPass        = "mainflux"
	defDBName        = "users"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defHTTPPort      = "8180"
	defServerCert    = ""
	defServerKey     = ""
	defJaegerURL     = ""

	defAuthnHTTPPort = "8989"
	defAuthnGRPCPort = "8181"
	defAuthnTimeout  = "1" // in seconds
	defAuthnTLS      = "false"
	defAuthnCACerts  = ""
	defAuthnURL      = "localhost:8181"

	defEmailLogLevel    = "debug"
	defEmailDriver      = "smtp"
	defEmailHost        = "localhost"
	defEmailPort        = "25"
	defEmailUsername    = "root"
	defEmailPassword    = ""
	defEmailFromAddress = ""
	defEmailFromName    = ""
	defEmailTemplate    = "email.tmpl"

	defTokenResetEndpoint = "/reset-request" // URL where user lands after click on the reset link from email

	defCacheURL  = "localhost:6379"
	defCachePass = ""
	defCacheDB   = "0"

	defNatsURL   = nats.DefaultURL
	defchannelID = ""

	envLogLevel      = "PD_USERS_LOG_LEVEL"
	envDBHost        = "PD_USERS_DB_HOST"
	envDBPort        = "PD_USERS_DB_PORT"
	envDBUser        = "PD_USERS_DB_USER"
	envDBPass        = "PD_USERS_DB_PASS"
	envDBName        = "PD_USERS_DB"
	envDBSSLMode     = "PD_USERS_DB_SSL_MODE"
	envDBSSLCert     = "PD_USERS_DB_SSL_CERT"
	envDBSSLKey      = "PD_USERS_DB_SSL_KEY"
	envDBSSLRootCert = "PD_USERS_DB_SSL_ROOT_CERT"
	envHTTPPort      = "PD_USERS_HTTP_PORT"
	envServerCert    = "PD_USERS_SERVER_CERT"
	envServerKey     = "PD_USERS_SERVER_KEY"
	envJaegerURL     = "PD_JAEGER_URL"

	envAuthnHTTPPort = "PD_AUTHN_HTTP_PORT"
	envAuthnGRPCPort = "PD_AUTHN_GRPC_PORT"
	envAuthnTimeout  = "PD_AUTHN_TIMEOUT"
	envAuthnTLS      = "PD_AUTHN_CLIENT_TLS"
	envAuthnCACerts  = "PD_AUTHN_CA_CERTS"
	envAuthnURL      = "PD_AUTHN_URL"

	envEmailDriver      = "PD_EMAIL_DRIVER"
	envEmailHost        = "PD_EMAIL_HOST"
	envEmailPort        = "PD_EMAIL_PORT"
	envEmailUsername    = "PD_EMAIL_USERNAME"
	envEmailPassword    = "PD_EMAIL_PASSWORD"
	envEmailFromAddress = "PD_EMAIL_FROM_ADDRESS"
	envEmailFromName    = "PD_EMAIL_FROM_NAME"
	envEmailLogLevel    = "PD_EMAIL_LOG_LEVEL"
	envEmailTemplate    = "PD_EMAIL_TEMPLATE"

	envTokenResetEndpoint = "PD_TOKEN_RESET_ENDPOINT"

	envCacheURL  = "PD_RULECHAIN_CACHE_URL"
	envCachePass = "PD_RULECHAIN_CACHE_PASS"
	envCacheDB   = "PD_RULECHAIN_CACHE_DB"
	envNatsURL   = "PD_NATS_URL"
	envchannelID = "PD_RULECHAIN_CHANNEL_ID"
)

type config struct {
	logLevel      string
	dbConfig      postgres.Config
	authnHTTPPort string
	authnGRPCPort string
	authnTimeout  time.Duration
	authnTLS      bool
	authnCACerts  string
	authnURL      string
	emailConf     email.Config
	httpPort      string
	serverCert    string
	serverKey     string
	jaegerURL     string
	resetURL      string
	cacheURL      string
	cachePass     string
	cacheDB       string
	NatsURL       string
	channelID     string
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	cacheClient := connectToRedis(cfg.cacheURL, cfg.cachePass, cfg.cacheDB, logger)

	authTracer, closer := initJaeger("auth", cfg.jaegerURL, logger)
	defer closer.Close()

	auth, close := connectToAuthn(cfg, authTracer, logger)
	if close != nil {
		defer close()
	}

	tracer, closer := initJaeger("rulechain", cfg.jaegerURL, logger)
	defer closer.Close()

	dbTracer, dbCloser := initJaeger("rulechain_db", cfg.jaegerURL, logger)
	defer dbCloser.Close()

	cacheTracer, cacheCloser := initJaeger("rulechain_cache", cfg.jaegerURL, logger)
	defer cacheCloser.Close()

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		logger.Error(fmt.Sprint("failed to connect to NATS: %s", err))
		os.Exit(1)
	}
	defer nc.Close()

	svc := newService(nc, cfg.channelID, db, cacheClient, dbTracer, cacheTracer, auth, cfg, logger)
	errs := make(chan error, 2)

	go startHTTPServer(tracer, svc, cfg.httpPort, cfg.serverCert, cfg.serverKey, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Users service terminated: %s", err))
}

func loadConfig() config {
	timeout, err := strconv.ParseInt(pandas.Env(envAuthnTimeout, defAuthnTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envAuthnTimeout, err.Error())
	}

	tls, err := strconv.ParseBool(pandas.Env(envAuthnTLS, defAuthnTLS))
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envAuthnTLS)
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

	emailConf := email.Config{
		Driver:      pandas.Env(envEmailDriver, defEmailDriver),
		FromAddress: pandas.Env(envEmailFromAddress, defEmailFromAddress),
		FromName:    pandas.Env(envEmailFromName, defEmailFromName),
		Host:        pandas.Env(envEmailHost, defEmailHost),
		Port:        pandas.Env(envEmailPort, defEmailPort),
		Username:    pandas.Env(envEmailUsername, defEmailUsername),
		Password:    pandas.Env(envEmailPassword, defEmailPassword),
		Template:    pandas.Env(envEmailTemplate, defEmailTemplate),
	}

	return config{
		logLevel:      pandas.Env(envLogLevel, defLogLevel),
		dbConfig:      dbConfig,
		authnHTTPPort: pandas.Env(envAuthnHTTPPort, defAuthnHTTPPort),
		authnGRPCPort: pandas.Env(envAuthnGRPCPort, defAuthnGRPCPort),
		authnURL:      pandas.Env(envAuthnURL, defAuthnURL),
		authnTimeout:  time.Duration(timeout) * time.Second,
		authnTLS:      tls,
		emailConf:     emailConf,
		httpPort:      pandas.Env(envHTTPPort, defHTTPPort),
		serverCert:    pandas.Env(envServerCert, defServerCert),
		serverKey:     pandas.Env(envServerKey, defServerKey),
		jaegerURL:     pandas.Env(envJaegerURL, defJaegerURL),
		resetURL:      pandas.Env(envTokenResetEndpoint, defTokenResetEndpoint),
		cacheURL:      pandas.Env(envCacheURL, defCacheURL),
		cachePass:     pandas.Env(envCachePass, defCachePass),
		cacheDB:       pandas.Env(envCacheDB, defCacheDB),
		NatsURL:       pandas.Env(envNatsURL, defNatsURL),
		channelID:     pandas.Env(envchannelID, defchannelID),
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
		logger.Error(fmt.Sprintf("Failed to init Jaeger: %s", err))
		os.Exit(1)
	}

	return tracer, closer
}

func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sqlx.DB {
	db, err := postgres.Connect(dbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to postgres: %s", err))
		os.Exit(1)
	}

	return db
}

func connectToRedis(cacheURL string, cachePass string, cacheDB string, logger logger.Logger) *redis.Client {
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

func connectToAuthn(cfg config, tracer opentracing.Tracer, logger logger.Logger) (mainflux.AuthNServiceClient, func() error) {
	var opts []grpc.DialOption
	if cfg.authnTLS {
		if cfg.authnCACerts != "" {
			tpc, err := credentials.NewClientTLSFromFile(cfg.authnCACerts, "")
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
		logger.Error(fmt.Sprintf("Failed to connect to users service: %s", err))
		os.Exit(1)
	}

	return authapi.NewClient(tracer, conn, cfg.authnTimeout), conn.Close
}

func newService(nc *nats.Conn, chanID string, db *sqlx.DB, cacheClient *redis.Client, dbTracer opentracing.Tracer, cacheTracer opentracing.Tracer, auth mainflux.AuthNServiceClient, c config, logger logger.Logger) rulechain.Service {
	database := postgres.NewDatabase(db)

	repo := tracing.RulechainRepositoryMiddleware(postgres.NewRuleChainRepository(database), dbTracer)

	cache := rediscache.NewRuleChainCache(cacheClient)
	cache = tracing.RuleChainCacheMiddleware(cacheTracer, cache)

	instancemanager := rulechain.NewInstanceManager()
	svc := rulechain.New(auth, repo, *instancemanager, cache)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "rulechain",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "rulechain",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	natssub.NewSubscriber(nc, chanID, svc, logger)
	return svc
}

func startHTTPServer(tracer opentracing.Tracer, svc rulechain.Service, port string, certFile string, keyFile string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if certFile != "" || keyFile != "" {
		logger.Info(fmt.Sprintf("Users service started using https, cert %s key %s, exposed port %s", certFile, keyFile, port))
		errs <- http.ListenAndServeTLS(p, certFile, keyFile, rulechainapi.MakeHandler(svc, tracer, logger))
	} else {
		logger.Info(fmt.Sprintf("Users service started using http, exposed port %s", port))
		errs <- http.ListenAndServe(p, rulechainapi.MakeHandler(svc, tracer, logger))
	}
}
