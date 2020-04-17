//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/lbs"
	"github.com/cloustone/pandas/lbs/api"
	lbshttpapi "github.com/cloustone/pandas/lbs/api/http"
	lbp "github.com/cloustone/pandas/lbs/proxy"
	"github.com/cloustone/pandas/pkg/logger"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	opentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
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

	envLogLevel        = "PD_THINGS_LOG_LEVEL"
	envDBHost          = "PD_THINGS_DB_HOST"
	envDBPort          = "PD_THINGS_DB_PORT"
	envDBUser          = "PD_THINGS_DB_USER"
	envDBPass          = "PD_THINGS_DB_PASS"
	envDBName          = "PD_THINGS_DB"
	envDBSSLMode       = "PD_THINGS_DB_SSL_MODE"
	envDBSSLCert       = "PD_THINGS_DB_SSL_CERT"
	envDBSSLKey        = "PD_THINGS_DB_SSL_KEY"
	envDBSSLRootCert   = "PD_THINGS_DB_SSL_ROOT_CERT"
	envClientTLS       = "PD_THINGS_CLIENT_TLS"
	envCACerts         = "PD_THINGS_CA_CERTS"
	envCacheURL        = "PD_THINGS_CACHE_URL"
	envCachePass       = "PD_THINGS_CACHE_PASS"
	envCacheDB         = "PD_THINGS_CACHE_DB"
	envESURL           = "PD_THINGS_ES_URL"
	envESPass          = "PD_THINGS_ES_PASS"
	envESDB            = "PD_THINGS_ES_DB"
	envHTTPPort        = "PD_THINGS_HTTP_PORT"
	envAuthHTTPPort    = "PD_THINGS_AUTH_HTTP_PORT"
	envAuthGRPCPort    = "PD_THINGS_AUTH_GRPC_PORT"
	envServerCert      = "PD_THINGS_SERVER_CERT"
	envServerKey       = "PD_THINGS_SERVER_KEY"
	envSingleUserEmail = "PD_THINGS_SINGLE_USER_EMAIL"
	envSingleUserToken = "PD_THINGS_SINGLE_USER_TOKEN"
	envJaegerURL       = "PD_JAEGER_URL"
	envAuthURL         = "PD_AUTH_URL"
	envAuthTimeout     = "PD_AUTH_TIMEOUT"
)

// inject by go build
var (
	Version   = "0.0.0"
	BuildTime = "2020-01-13-0802 UTC"
)

type config struct {
	logLevel        string
	clientTLS       bool
	caCerts         string
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

func init() {
	fmt.Println("Version:", Version)
	fmt.Println("BuildTime:", BuildTime)
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	lbsTracer, lbsCloser := initJaeger("lbs", cfg.jaegerURL, logger)
	defer lbsCloser.Close()

	location := lbp.NewLocationServingOptions()
	svc := newService(location, logger)
	errs := make(chan error, 2)

	go startHTTPServer(lbshttpapi.MakeHandler(lbsTracer, svc), cfg.httpPort, cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Lbs service terminated: %s", err))
	rand.Seed(time.Now().UTC().UnixNano())

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
		logLevel:  pandas.Env(envLogLevel, defLogLevel),
		clientTLS: tls,
		caCerts:   pandas.Env(envCACerts, defCACerts),

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

func newService(location *lbp.LocationServingOptions, logger logger.Logger) lbs.Service {

	proxy := lbp.NewProxy(location)

	svc := lbs.New(*proxy)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "lbs",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "lbs",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
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
func startHTTPServer(handler http.Handler, port string, cfg config, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if cfg.serverCert != "" || cfg.serverKey != "" {
		logger.Info(fmt.Sprintf("Things service started using https on port %s with cert %s key %s",
			port, cfg.serverCert, cfg.serverKey))
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, handler)
		return
	}
	logger.Info(fmt.Sprintf("Things service started using http on port %s", cfg.httpPort))
	errs <- http.ListenAndServe(p, handler)
}
