package nats

import (
	"encoding/json"
	"log"

	"github.com/cloustone/pandas/pkg/broadcast"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const NAME = "nats"

type NatsBroadcast struct {
	nc          *nats.Conn
	subscribers map[string]subscriber
}

type subscriber struct {
	observer broadcast.Observer
}

func NewBroadcast(natsurl string) broadcast.Broadcast {
	//nats.DefaultURL
	natsconn, err := nats.Connect(natsurl)
	if err != nil {
		log.Fatal(err)
	}
	return &NatsBroadcast{
		nc:          natsconn,
		subscribers: make(map[string]subscriber),
	}
}

func (n *NatsBroadcast) AsMember() {}

func (n *NatsBroadcast) WithRootPath(path string) broadcast.Broadcast { return n }

func (n *NatsBroadcast) Notify(no broadcast.Notification) {
	body, err := json.Marshal(&no)
	if err != nil {
		logrus.WithError(err)
	}
	if err = n.nc.Publish(no.ObjectPath, body); err != nil {
		logrus.WithError(err)
	}
}

func (n *NatsBroadcast) RegisterObserver(path string, obs broadcast.Observer) {
	subscri := subscriber{
		observer: obs,
	}
	go async(n, n.nc, path, subscri)
	n.subscribers[path] = subscri
}

func async(n *NatsBroadcast, nc *nats.Conn, path string, sub subscriber) {
	nc.Subscribe(path, func(msg *nats.Msg) {
		no := broadcast.Notification{}
		if err := json.Unmarshal(msg.Data, &no); err != nil {
			logrus.WithError(err)
		}
		sub.observer.Onbroadcast(n, no)
	})
}
