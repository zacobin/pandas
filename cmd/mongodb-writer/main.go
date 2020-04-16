// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/mainflux/broker"
	"github.com/cloustone/pandas/mainflux/transformers/senml"
	"github.com/cloustone/pandas/mainflux/writers"
	"github.com/cloustone/pandas/mainflux/writers/api"
	"github.com/cloustone/pandas/mainflux/writers/mongodb"
	"github.com/cloustone/pandas/pkg/logger"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	svcName = "mongodb-writer"

	defNatsURL         = pandas.DefNatsURL
	defLogLevel        = "error"
	defPort            = "8180"
	defDBName          = "mainflux"
	defDBHost          = "localhost"
	defDBPort          = "27017"
	defSubjectsCfgPath = "/config/subjects.toml"

	envNatsURL         = "PD_NATS_URL"
	envLogLevel        = "PD_MONGO_WRITER_LOG_LEVEL"
	envPort            = "PD_MONGO_WRITER_PORT"
	envDBName          = "PD_MONGO_WRITER_DB_NAME"
	envDBHost          = "PD_MONGO_WRITER_DB_HOST"
	envDBPort          = "PD_MONGO_WRITER_DB_PORT"
	envSubjectsCfgPath = "PD_MONGO_WRITER_SUBJECTS_CONFIG"
)

type config struct {
	natsURL         string
	logLevel        string
	port            string
	dbName          string
	dbHost          string
	dbPort          string
	subjectsCfgPath string
}

func main() {
	cfg := loadConfigs()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatal(err)
	}

	b, err := broker.New(cfg.natsURL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer b.Close()

	addr := fmt.Sprintf("mongodb://%s:%s", cfg.dbHost, cfg.dbPort)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(addr))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to database: %s", err))
		os.Exit(1)
	}

	db := client.Database(cfg.dbName)
	repo := mongodb.New(db)

	counter, latency := makeMetrics()
	repo = api.LoggingMiddleware(repo, logger)
	repo = api.MetricsMiddleware(repo, counter, latency)
	st := senml.New()
	if err := writers.Start(b, repo, st, svcName, cfg.subjectsCfgPath, logger); err != nil {
		logger.Error(fmt.Sprintf("Failed to start MongoDB writer: %s", err))
		os.Exit(1)
	}

	errs := make(chan error, 2)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go startHTTPService(cfg.port, logger, errs)

	err = <-errs
	logger.Error(fmt.Sprintf("MongoDB writer service terminated: %s", err))
}

func loadConfigs() config {
	return config{
		natsURL:         pandas.Env(envNatsURL, defNatsURL),
		logLevel:        pandas.Env(envLogLevel, defLogLevel),
		port:            pandas.Env(envPort, defPort),
		dbName:          pandas.Env(envDBName, defDBName),
		dbHost:          pandas.Env(envDBHost, defDBHost),
		dbPort:          pandas.Env(envDBPort, defDBPort),
		subjectsCfgPath: pandas.Env(envSubjectsCfgPath, defSubjectsCfgPath),
	}
}

func makeMetrics() (*kitprometheus.Counter, *kitprometheus.Summary) {
	counter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "mongodb",
		Subsystem: "message_writer",
		Name:      "request_count",
		Help:      "Number of database inserts.",
	}, []string{"method"})

	latency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "mongodb",
		Subsystem: "message_writer",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of inserts in microseconds.",
	}, []string{"method"})

	return counter, latency
}

func startHTTPService(port string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Mongodb writer service started, exposed port %s", p))
	errs <- http.ListenAndServe(p, api.MakeHandler(svcName))
}
