package main

import (
	"context"
	"fmt"

	"github.com/cloustone/pandas/kuiper/xstream/api"
)

type natsSource struct {
	subscriber *nats.Subscriber
	srv        string
	topic      string
	cancel     context.CancelFunc
}

func (s *natsSource) Configure(topic string, props map[string]interface{}) error {
	s.topic = topic
	srv, ok := props["server"]
	if !ok {
		return fmt.Errorf("nats source is missing property server")
	}
	s.srv = srv.(string)
	return nil
}

func (s *natsSource) Open(ctx api.StreamContext, consumer chan<- api.SourceTuple, errCh chan<- error) {
}

func (s *natsSource) Close(ctx api.StreamContext) error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

func Nats() api.Source {
	return &natsSource{}
}
