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
	"github.com/cloustone/pandas/mainflux/readers/influxdb"
	"github.com/cloustone/pandas/pkg/logger"
	thingsapi "github.com/cloustone/pandas/things/api/auth/grpc"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	influxdata "github.com/influxdata/influxdb/client/v2"
	opentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defThingsURL     = "localhost:8181"
	defLogLevel      = "error"
	defPort          = "8180"
	defDBName        = "mainflux"
	defDBHost        = "localhost"
	defDBPort        = "8086"
	defDBUser        = "mainflux"
	defDBPass        = "mainflux"
	defClientTLS     = "false"
	defCACerts       = ""
	defServerCert    = ""
	defServerKey     = ""
	defJaegerURL     = ""
	defThingsTimeout = "1" // in seconds

	envThingsURL     = "MF_THINGS_URL"
	envLogLevel      = "MF_INFLUX_READER_LOG_LEVEL"
	envPort          = "MF_INFLUX_READER_PORT"
	envDBName        = "MF_INFLUX_READER_DB_NAME"
	envDBHost        = "MF_INFLUX_READER_DB_HOST"
	envDBPort        = "MF_INFLUX_READER_DB_PORT"
	envDBUser        = "MF_INFLUX_READER_DB_USER"
	envDBPass        = "MF_INFLUX_READER_DB_PASS"
	envClientTLS     = "MF_INFLUX_READER_CLIENT_TLS"
	envCACerts       = "MF_INFLUX_READER_CA_CERTS"
	envServerCert    = "MF_INFLUX_READER_SERVER_CERT"
	envServerKey     = "MF_INFLUX_READER_SERVER_KEY"
	envJaegerURL     = "MF_JAEGER_URL"
	envThingsTimeout = "MF_INFLUX_READER_THINGS_TIMEOUT"
)

type config struct {
	thingsURL     string
	logLevel      string
	port          string
	dbName        string
	dbHost        string
	dbPort        string
	dbUser        string
	dbPass        string
	clientTLS     bool
	caCerts       string
	serverCert    string
	serverKey     string
	jaegerURL     string
	thingsTimeout time.Duration
}

func main() {
	cfg, clientCfg := loadConfigs()
	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}
	conn := connectToThings(cfg, logger)
	defer conn.Close()

	thingsTracer, thingsCloser := initJaeger("things", cfg.jaegerURL, logger)
	defer thingsCloser.Close()

	tc := thingsapi.NewClient(conn, thingsTracer, cfg.thingsTimeout)

	client, err := influxdata.NewHTTPClient(clientCfg)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create InfluxDB client: %s", err))
		os.Exit(1)
	}
	defer client.Close()

	repo := newService(client, cfg.dbName, logger)

	errs := make(chan error, 2)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go startHTTPServer(repo, tc, cfg, logger, errs)

	err = <-errs
	logger.Error(fmt.Sprintf("InfluxDB writer service terminated: %s", err))
}

func loadConfigs() (config, influxdata.HTTPConfig) {
	tls, err := strconv.ParseBool(pandas.Env(envClientTLS, defClientTLS))
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envClientTLS)
	}

	timeout, err := strconv.ParseInt(pandas.Env(envThingsTimeout, defThingsTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envThingsTimeout, err.Error())
	}

	cfg := config{
		thingsURL:     pandas.Env(envThingsURL, defThingsURL),
		logLevel:      pandas.Env(envLogLevel, defLogLevel),
		port:          pandas.Env(envPort, defPort),
		dbName:        pandas.Env(envDBName, defDBName),
		dbHost:        pandas.Env(envDBHost, defDBHost),
		dbPort:        pandas.Env(envDBPort, defDBPort),
		dbUser:        pandas.Env(envDBUser, defDBUser),
		dbPass:        pandas.Env(envDBPass, defDBPass),
		clientTLS:     tls,
		caCerts:       pandas.Env(envCACerts, defCACerts),
		serverCert:    pandas.Env(envServerCert, defServerCert),
		serverKey:     pandas.Env(envServerKey, defServerKey),
		jaegerURL:     pandas.Env(envJaegerURL, defJaegerURL),
		thingsTimeout: time.Duration(timeout) * time.Second,
	}

	clientCfg := influxdata.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%s", cfg.dbHost, cfg.dbPort),
		Username: cfg.dbUser,
		Password: cfg.dbPass,
	}

	return cfg, clientCfg
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

func newService(client influxdata.Client, dbName string, logger logger.Logger) readers.MessageRepository {
	repo := influxdb.New(client, dbName)
	repo = api.LoggingMiddleware(repo, logger)
	repo = api.MetricsMiddleware(
		repo,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "influxdb",
			Subsystem: "message_reader",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "influxdb",
			Subsystem: "message_reader",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return repo
}

func startHTTPServer(repo readers.MessageRepository, tc mainflux.ThingsServiceClient, cfg config, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.port)
	if cfg.serverCert != "" || cfg.serverKey != "" {
		logger.Info(fmt.Sprintf("InfluxDB reader service started using https on port %s with cert %s key %s",
			cfg.port, cfg.serverCert, cfg.serverKey))
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, api.MakeHandler(repo, tc, "influxdb-reader"))
		return
	}
	logger.Info(fmt.Sprintf("InfluxDB reader service started, exposed port %s", cfg.port))
	errs <- http.ListenAndServe(p, api.MakeHandler(repo, tc, "influxdb-reader"))
}
