// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
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

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/email"
	"github.com/cloustone/pandas/realms"
	"github.com/cloustone/pandas/realms/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	authapi "github.com/cloustone/pandas/authn/api/grpc"
	"github.com/cloustone/pandas/pkg/logger"
	"github.com/cloustone/pandas/realms/api"
	"github.com/cloustone/pandas/realms/bcrypt"
	"github.com/cloustone/pandas/realms/postgres"
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
	defDBName        = "realms"
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

	defRealmsFile = "./realms.json"

	defTokenResetEndpoint = "/reset-request" // URL where user lands after click on the reset link from email

	envLogLevel      = "MF_USERS_LOG_LEVEL"
	envDBHost        = "MF_USERS_DB_HOST"
	envDBPort        = "MF_USERS_DB_PORT"
	envDBUser        = "MF_USERS_DB_USER"
	envDBPass        = "MF_USERS_DB_PASS"
	envDBName        = "MF_USERS_DB"
	envDBSSLMode     = "MF_USERS_DB_SSL_MODE"
	envDBSSLCert     = "MF_USERS_DB_SSL_CERT"
	envDBSSLKey      = "MF_USERS_DB_SSL_KEY"
	envDBSSLRootCert = "MF_USERS_DB_SSL_ROOT_CERT"
	envHTTPPort      = "MF_USERS_HTTP_PORT"
	envServerCert    = "MF_USERS_SERVER_CERT"
	envServerKey     = "MF_USERS_SERVER_KEY"
	envJaegerURL     = "MF_JAEGER_URL"

	envAuthnHTTPPort = "MF_AUTHN_HTTP_PORT"
	envAuthnGRPCPort = "MF_AUTHN_GRPC_PORT"
	envAuthnTimeout  = "MF_AUTHN_TIMEOUT"
	envAuthnTLS      = "MF_AUTHN_CLIENT_TLS"
	envAuthnCACerts  = "MF_AUTHN_CA_CERTS"
	envAuthnURL      = "MF_AUTHN_URL"

	envEmailDriver      = "MF_EMAIL_DRIVER"
	envEmailHost        = "MF_EMAIL_HOST"
	envEmailPort        = "MF_EMAIL_PORT"
	envEmailUsername    = "MF_EMAIL_USERNAME"
	envEmailPassword    = "MF_EMAIL_PASSWORD"
	envEmailFromAddress = "MF_EMAIL_FROM_ADDRESS"
	envEmailFromName    = "MF_EMAIL_FROM_NAME"
	envEmailLogLevel    = "MF_EMAIL_LOG_LEVEL"
	envEmailTemplate    = "MF_EMAIL_TEMPLATE"

	envRealmsFile = "MF_REAMS_FILE"

	envTokenResetEndpoint = "MF_TOKEN_RESET_ENDPOINT"
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
	realmsFile    string
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	authTracer, closer := initJaeger("auth", cfg.jaegerURL, logger)
	defer closer.Close()

	auth, close := connectToAuthn(cfg, authTracer, logger)
	if close != nil {
		defer close()
	}

	tracer, closer := initJaeger("realms", cfg.jaegerURL, logger)
	defer closer.Close()

	dbTracer, dbCloser := initJaeger("realms_db", cfg.jaegerURL, logger)
	defer dbCloser.Close()

	svc := newService(db, dbTracer, auth, cfg, logger)
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
	timeout, err := strconv.ParseInt(mainflux.Env(envAuthnTimeout, defAuthnTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envAuthnTimeout, err.Error())
	}

	tls, err := strconv.ParseBool(mainflux.Env(envAuthnTLS, defAuthnTLS))
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envAuthnTLS)
	}

	dbConfig := postgres.Config{
		Host:        mainflux.Env(envDBHost, defDBHost),
		Port:        mainflux.Env(envDBPort, defDBPort),
		User:        mainflux.Env(envDBUser, defDBUser),
		Pass:        mainflux.Env(envDBPass, defDBPass),
		Name:        mainflux.Env(envDBName, defDBName),
		SSLMode:     mainflux.Env(envDBSSLMode, defDBSSLMode),
		SSLCert:     mainflux.Env(envDBSSLCert, defDBSSLCert),
		SSLKey:      mainflux.Env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: mainflux.Env(envDBSSLRootCert, defDBSSLRootCert),
	}

	emailConf := email.Config{
		Driver:      mainflux.Env(envEmailDriver, defEmailDriver),
		FromAddress: mainflux.Env(envEmailFromAddress, defEmailFromAddress),
		FromName:    mainflux.Env(envEmailFromName, defEmailFromName),
		Host:        mainflux.Env(envEmailHost, defEmailHost),
		Port:        mainflux.Env(envEmailPort, defEmailPort),
		Username:    mainflux.Env(envEmailUsername, defEmailUsername),
		Password:    mainflux.Env(envEmailPassword, defEmailPassword),
		Template:    mainflux.Env(envEmailTemplate, defEmailTemplate),
	}

	return config{
		logLevel:      mainflux.Env(envLogLevel, defLogLevel),
		dbConfig:      dbConfig,
		authnHTTPPort: mainflux.Env(envAuthnHTTPPort, defAuthnHTTPPort),
		authnGRPCPort: mainflux.Env(envAuthnGRPCPort, defAuthnGRPCPort),
		authnURL:      mainflux.Env(envAuthnURL, defAuthnURL),
		authnTimeout:  time.Duration(timeout) * time.Second,
		authnTLS:      tls,
		emailConf:     emailConf,
		httpPort:      mainflux.Env(envHTTPPort, defHTTPPort),
		serverCert:    mainflux.Env(envServerCert, defServerCert),
		serverKey:     mainflux.Env(envServerKey, defServerKey),
		jaegerURL:     mainflux.Env(envJaegerURL, defJaegerURL),
		resetURL:      mainflux.Env(envTokenResetEndpoint, defTokenResetEndpoint),
		realmsFile:    mainflux.Env(envRealmsFile, defRealmsFile),
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
		logger.Error(fmt.Sprintf("Failed to connect to realms service: %s", err))
		os.Exit(1)
	}

	return authapi.NewClient(tracer, conn, cfg.authnTimeout), conn.Close
}

func loadRealmProviders(c config, repo realms.RealmRepository, logger logger.Logger) ([]realms.RealmProvider, error) {
	realmList := []realms.Realm{}
	providers := []realms.RealmProvider{}

	// Load realms from repository, if there are no realms, then load from
	// config file and push these realms into repository
	if rs, err := repo.List(context.TODO()); err != nil {
		logger.Debug(fmt.Sprintf("There are no realms in repository :%s", err.Error()))
	} else {
		realmList = append(realmList, rs...)
	}
	// Create realm provider from builtin realm config file
	if rs, err := loadRealmsWithFile(c.realmsFile, logger); err != nil {
		logger.Debug(fmt.Sprintf("Failed to load default realms config file: %s", err.Error()))
	} else {
		if len(realmList) == 0 { // no realms already existed in repository
			for _, realm := range rs {
				repo.Save(context.TODO(), realm)
			}
		}
		realmList = append(realmList, rs...)
	}

	// All realms had been loaded
	for _, realm := range realmList {
		provider, err := realms.NewRealmProvider(realm)
		if err != nil {
			logger.Debug(fmt.Sprintf("Failed to create realm provider for '%s'", realm.Name))
			continue
		}
		providers = append(providers, provider)
	}
	if len(providers) == 0 {
		logger.Debug(fmt.Sprintf("No realm provider already exist"))
	}
	return providers, nil
}

// loadReamsWithFile load realms from realms config file
func loadRealmsWithFile(fullFilePath string, logger logger.Logger) ([]realms.Realm, error) {
	buf, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		logger.Debug("open realms config file failed")
		return nil, err
	}
	realms := []realms.Realm{}
	if err := json.Unmarshal(buf, &realms); err != nil {
		logger.Debug("illegal realm config file")
		return nil, err
	}
	return realms, nil
}

func newService(db *sqlx.DB, tracer opentracing.Tracer, auth mainflux.AuthNServiceClient, c config, logger logger.Logger) realms.Service {
	database := postgres.NewDatabase(db)
	repo := tracing.RealmRepositoryMiddleware(postgres.New(database), tracer)
	hasher := bcrypt.New()
	providers, err := loadRealmProviders(c, repo, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load realm providers: %s", err.Error()))
	}

	svc := realms.New(repo, hasher, auth, providers)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "realms",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "realms",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}

func startHTTPServer(tracer opentracing.Tracer, svc realms.Service, port string, certFile string, keyFile string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if certFile != "" || keyFile != "" {
		logger.Info(fmt.Sprintf("Realms service started using https, cert %s key %s, exposed port %s", certFile, keyFile, port))
		errs <- http.ListenAndServeTLS(p, certFile, keyFile, api.MakeHandler(svc, tracer, logger))
	} else {
		logger.Info(fmt.Sprintf("Realms service started using http, exposed port %s", port))
		errs <- http.ListenAndServe(p, api.MakeHandler(svc, tracer, logger))
	}
}
