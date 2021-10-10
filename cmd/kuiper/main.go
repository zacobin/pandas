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
	"github.com/cloustone/pandas/kuiper/plugins"
	"github.com/cloustone/pandas/kuiper/tracing"
	"github.com/cloustone/pandas/mainflux"

	"github.com/jmoiron/sqlx"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/credentials"

	authapi "github.com/cloustone/pandas/authn/api/grpc"
	"github.com/cloustone/pandas/kuiper"
	"github.com/cloustone/pandas/kuiper/api"
	thhttpapi "github.com/cloustone/pandas/kuiper/api/http"
	"github.com/cloustone/pandas/kuiper/postgres"
	rediscache "github.com/cloustone/pandas/kuiper/redis"
	"github.com/cloustone/pandas/kuiper/uuid"
	"github.com/cloustone/pandas/pkg/logger"
	localusers "github.com/cloustone/pandas/things/users"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-redis/redis"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
)

const (
	defLogLevel        = "info"
	defDBHost          = "localhost"
	defDBPort          = "5432"
	defDBUser          = "pandas"
	defDBPass          = "grgecent"
	defDBName          = "kuiper"
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

	envLogLevel        = "PD_KUIPER_LOG_LEVEL"
	envDBHost          = "PD_KUIPER_DB_HOST"
	envDBPort          = "PD_KUIPER_DB_PORT"
	envDBUser          = "PD_KUIPER_DB_USER"
	envDBPass          = "PD_KUIPER_DB_PASS"
	envDBName          = "PD_KUIPER_DB"
	envDBSSLMode       = "PD_KUIPER_DB_SSL_MODE"
	envDBSSLCert       = "PD_KUIPER_DB_SSL_CERT"
	envDBSSLKey        = "PD_KUIPER_DB_SSL_KEY"
	envDBSSLRootCert   = "PD_KUIPER_DB_SSL_ROOT_CERT"
	envClientTLS       = "PD_KUIPER_CLIENT_TLS"
	envCACerts         = "PD_KUIPER_CA_CERTS"
	envCacheURL        = "PD_KUIPER_CACHE_URL"
	envCachePass       = "PD_KUIPER_CACHE_PASS"
	envCacheDB         = "PD_KUIPER_CACHE_DB"
	envESURL           = "PD_KUIPER_ES_URL"
	envESPass          = "PD_KUIPER_ES_PASS"
	envESDB            = "PD_KUIPER_ES_DB"
	envHTTPPort        = "PD_KUIPER_HTTP_PORT"
	envAuthHTTPPort    = "PD_KUIPER_AUTH_HTTP_PORT"
	envAuthGRPCPort    = "PD_KUIPER_AUTH_GRPC_PORT"
	envServerCert      = "PD_KUIPER_SERVER_CERT"
	envServerKey       = "PD_KUIPER_SERVER_KEY"
	envSingleUserEmail = "PD_KUIPER_SINGLE_USER_EMAIL"
	envSingleUserToken = "PD_KUIPER_SINGLE_USER_TOKEN"
	envJaegerURL       = "PD_JAEGER_URL"
	envAuthURL         = "PD_AUTH_URL"
	envAuthTimeout     = "PD_AUTH_TIMEOUT"
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
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	kuiperTracer, kuiperCloser := initJaeger("kuiper", cfg.jaegerURL, logger)
	defer kuiperCloser.Close()

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

	dbTracer, dbCloser := initJaeger("kuiper_db", cfg.jaegerURL, logger)
	defer dbCloser.Close()

	cacheTracer, cacheCloser := initJaeger("kuiper_cache", cfg.jaegerURL, logger)
	defer cacheCloser.Close()

	svc := newService(auth, dbTracer, cacheTracer, db, cacheClient, esClient, logger)
	errs := make(chan error, 2)

	go startHTTPServer(thhttpapi.MakeHandler(kuiperTracer, svc), cfg.httpPort, cfg, logger, errs)
	//	go startHTTPServer(authhttpapi.MakeHandler(kuiperTracer, svc), cfg.authHTTPPort, cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Kuiper service terminated: %s", err))
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

func newService(auth mainflux.AuthNServiceClient, dbTracer opentracing.Tracer, cacheTracer opentracing.Tracer, db *sqlx.DB, cacheClient *redis.Client, esClient *redis.Client, logger logger.Logger) kuiper.Service {
	database := postgres.NewDatabase(db)

	ruleRepo := postgres.NewRuleRepository(database)
	ruleRepo = tracing.RuleRepositoryMiddleware(dbTracer, ruleRepo)

	streamRepo := postgres.NewStreamRepository(database)
	streamRepo = tracing.StreamRepositoryMiddleware(dbTracer, streamRepo)

	ruleCache := rediscache.NewRuleCache(cacheClient)
	ruleCache = tracing.RuleCacheMiddleware(cacheTracer, ruleCache)

	streamCache := rediscache.NewStreamCache(cacheClient)
	streamCache = tracing.StreamCacheMiddleware(cacheTracer, streamCache)
	idp := uuid.New()

	pluginManager, err := plugins.NewPluginManager()
	if err != nil {
		log.Fatalf(err.Error())
	}

	svc := kuiper.New(auth, streamRepo, ruleRepo, streamCache, ruleCache, idp, pluginManager)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "kuiper",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "kuiper",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func startHTTPServer(handler http.Handler, port string, cfg config, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if cfg.serverCert != "" || cfg.serverKey != "" {
		logger.Info(fmt.Sprintf("Kuiper service started using https on port %s with cert %s key %s",
			port, cfg.serverCert, cfg.serverKey))
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, handler)
		return
	}
	logger.Info(fmt.Sprintf("Kuiper service started using http on port %s", cfg.httpPort))
	errs <- http.ListenAndServe(p, handler)
}
