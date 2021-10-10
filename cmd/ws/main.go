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
	"github.com/cloustone/pandas/mainflux/broker"
	adapter "github.com/cloustone/pandas/mainflux/ws"
	"github.com/cloustone/pandas/mainflux/ws/api"
	"github.com/cloustone/pandas/pkg/logger"
	thingsapi "github.com/cloustone/pandas/things/api/auth/grpc"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	opentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defClientTLS     = "false"
	defCACerts       = ""
	defPort          = "8180"
	defLogLevel      = "error"
	defNatsURL       = pandas.DefNatsURL
	defThingsURL     = "localhost:8181"
	defJaegerURL     = ""
	defThingsTimeout = "1" // in seconds

	envClientTLS     = "PD_WS_ADAPTER_CLIENT_TLS"
	envCACerts       = "PD_WS_ADAPTER_CA_CERTS"
	envPort          = "PD_WS_ADAPTER_PORT"
	envLogLevel      = "PD_WS_ADAPTER_LOG_LEVEL"
	envNatsURL       = "PD_NATS_URL"
	envThingsURL     = "PD_THINGS_URL"
	envJaegerURL     = "PD_JAEGER_URL"
	envThingsTimeout = "PD_WS_ADAPTER_THINGS_TIMEOUT"
)

type config struct {
	clientTLS     bool
	caCerts       string
	thingsURL     string
	natsURL       string
	logLevel      string
	port          string
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

	cc := thingsapi.NewClient(conn, thingsTracer, cfg.thingsTimeout)

	b, err := broker.New(cfg.natsURL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer b.Close()

	svc := newService(b, logger)

	errs := make(chan error, 2)

	go func() {
		p := fmt.Sprintf(":%s", cfg.port)
		logger.Info(fmt.Sprintf("WebSocket adapter service started, exposed port %s", cfg.port))
		errs <- http.ListenAndServe(p, api.MakeHandler(svc, cc, logger))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("WebSocket adapter terminated: %s", err))
}

func loadConfig() config {
	tls, err := strconv.ParseBool(pandas.Env(envClientTLS, defClientTLS))
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envClientTLS)
	}

	timeout, err := strconv.ParseInt(pandas.Env(envThingsTimeout, defThingsTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envThingsTimeout, err.Error())
	}

	return config{
		clientTLS:     tls,
		caCerts:       pandas.Env(envCACerts, defCACerts),
		thingsURL:     pandas.Env(envThingsURL, defThingsURL),
		natsURL:       pandas.Env(envNatsURL, defNatsURL),
		logLevel:      pandas.Env(envLogLevel, defLogLevel),
		port:          pandas.Env(envPort, defPort),
		jaegerURL:     pandas.Env(envJaegerURL, defJaegerURL),
		thingsTimeout: time.Duration(timeout) * time.Second,
	}
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

func newService(broker broker.Nats, log logger.Logger) adapter.Service {
	svc := adapter.New(broker, log)
	svc = api.LoggingMiddleware(svc, log)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "ws_adapter",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "ws_adapter",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}
