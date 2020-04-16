// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package subscriber

import (
	"fmt"
	"os"

	"github.com/cloustone/pandas/mainflux"
	log "github.com/cloustone/pandas/pkg/logger"
	"github.com/cloustone/pandas/rulechain"
	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/nats.go"
)

const (
	subject    = "channel.>"
	queuegroup = "rulechain"
)

//Subscriber subscriber
type Subscriber struct {
	natsClient *nats.Conn
	logger     log.Logger
	svc        rulechain.Service
	channelID  string
}

//NewSubscriber newsubscriber
func NewSubscriber(nc *nats.Conn, chID string, svc rulechain.Service, logger log.Logger) *Subscriber {
	s := Subscriber{
		natsClient: nc,
		logger:     logger,
		svc:        svc,
		channelID:  chID,
	}
	if _, err := s.natsClient.QueueSubscribe(subject, queuegroup, s.handleMsg); err != nil {
		logger.Error(fmt.Sprint("Failed to subscribe to NATS: %s", err))
		os.Exit(1)
	}
	return &s
}

func (s *Subscriber) handleMsg(m *nats.Msg) {
	var msg mainflux.Message
	if err := proto.Unmarshal(m.Data, &msg); err != nil {
		s.logger.Warn(fmt.Sprintf("Unmarshalling failed: %s", err))
		return
	}

	if msg.Channel == s.channelID {
		return
	}

	if err := s.svc.SaveStates(&msg); err != nil {
		s.logger.Error(fmt.Sprintf("State save failed: %s", err))
		return
	}
}
