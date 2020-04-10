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
package nodes

import (
	"fmt"
	"time"

	"github.com/cloustone/pandas/apimachinery/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type externalMqttNode struct {
	bareNode
	TopicPattern      string
	Host              string
	Port              string
	ConnectTimeoutSec int
	ClientId          string
	CleanSession      bool
	Ssl               bool
	MqttCli           mqtt.Client
}

type externalMqttNodeFactory struct{}

func (f externalMqttNodeFactory) Name() string     { return "ExternalMqttNode" }
func (f externalMqttNodeFactory) Category() string { return NODE_CATEGORY_EXTERNAL }
func (f externalMqttNodeFactory) Create(id string, meta Metadata) (Node, error) {
	labels := []string{"Success", "Failure"}

	node := &externalMqttNode{
		bareNode: newBareNode(f.Name(), id, meta, labels),
	}
	broker := fmt.Sprintf("tcp://%s:%s", node.Host, node.Port)
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(node.ClientId)
	opts.SetCleanSession(node.CleanSession)
	opts.SetConnectTimeout(time.Duration(node.ConnectTimeoutSec) * time.Second)
	node.MqttCli = mqtt.NewClient(opts)

	if token := node.MqttCli.Connect(); token.Wait() && token.Error() != nil {
		logrus.WithError(token.Error())
		return nil, token.Error()
	}
	return decodePath(meta, node)
}

func (n *externalMqttNode) Handle(msg models.Message) error {
	logrus.Infof("%s handle message '%s'", n.Name(), msg.GetType())

	successLabelNode := n.GetLinkedNode("Success")
	failureLabelNode := n.GetLinkedNode("Failure")
	topic := n.TopicPattern //need fix add msg.metadata in it
	sendmqttmsg, err := msg.MarshalBinary()
	if err != nil {
		return nil
	}
	token := n.MqttCli.Publish(topic, 1, false, sendmqttmsg)
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return failureLabelNode.Handle(msg)
	}
	return successLabelNode.Handle(msg)
}
