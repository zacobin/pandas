package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/dashboard"
	"github.com/cloustone/pandas/pms/postgres"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	"github.com/cloustone/pandas/pkg/logger"
)

const (
	defLogLevel        = "info"
	defClientTLS       = "false"
	defCACerts         = ""
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
	envLogLevel        = "PD_DASHBOARD_LOG_LEVEL"
	envClientTLS       = "PD_DASHBOARD_CLIENT_TLS"
	envCACerts         = "PD_DASHBOARD_CA_CERTS"
	envHTTPPort        = "PD_DASHBOARD_HTTP_PORT"
	envAuthHTTPPort    = "PD_DASHBOARD_AUTH_HTTP_PORT"
	envAuthGRPCPort    = "PD_DASHBOARD_AUTH_GRPC_PORT"
	envServerCert      = "PD_DASHBOARD_SERVER_CERT"
	envServerKey       = "PD_DASHBOARD_SERVER_KEY"
	envSingleUserEmail = "PD_DASHBOARD_SINGLE_USER_EMAIL"
	envSingleUserToken = "PD_DASHBOARD_SINGLE_USER_TOKEN"
	envJaegerURL       = "PD_JAEGER_URL"
	envAuthURL         = "PD_AUTH_URL"
	envAuthTimeout     = "PD_AUTH_TIMEOUT"
	envDownstreams     = "PD_DOWNSTREAMS"
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
	downstreams     string
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	errs := make(chan error, 2)
	go startHTTPServer(cfg.httpPort, cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("dashboard service terminated: %s", err))
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

	return config{
		logLevel:        pandas.Env(envLogLevel, defLogLevel),
		clientTLS:       tls,
		caCerts:         pandas.Env(envCACerts, defCACerts),
		httpPort:        pandas.Env(envHTTPPort, defHTTPPort),
		authHTTPPort:    pandas.Env(envAuthHTTPPort, defAuthHTTPPort),
		serverCert:      pandas.Env(envServerCert, defServerCert),
		serverKey:       pandas.Env(envServerKey, defServerKey),
		singleUserEmail: pandas.Env(envSingleUserEmail, defSingleUserEmail),
		singleUserToken: pandas.Env(envSingleUserToken, defSingleUserToken),
		jaegerURL:       pandas.Env(envJaegerURL, defJaegerURL),
		authURL:         pandas.Env(envAuthURL, defAuthURL),
		authTimeout:     time.Duration(timeout) * time.Second,
	}
}

func startHTTPServer(port string, cfg config, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("dashboard service started using http on port %s", cfg.httpPort))
	http.Handle("/", http.FileServer(assetFS()))
	errs <- http.ListenAndServe(p, nil)
}

func assetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     dashboard.Asset,
		AssetDir:  dashboard.AssetDir,
		AssetInfo: dashboard.AssetInfo,
		Prefix:    "dist",
		Fallback:  "index.html",
	}
}
