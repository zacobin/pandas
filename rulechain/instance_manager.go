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
package rulechain

import (
	"fmt"
	"sync"

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/rulechain/message"
	logr "github.com/sirupsen/logrus"
)

// instanceManager manage all rulechain's runtime
type instanceManager struct {
	mutex      sync.RWMutex
	rulechains map[string]*ruleChainInstance
}

// newInstanceManager create controller instance used in rule chain service
func NewInstanceManager() *instanceManager {
	controller := &instanceManager{
		mutex:      sync.RWMutex{},
		rulechains: make(map[string]*ruleChainInstance),
	}
	return controller
}

// startRuleChain start the rule chain and receiving incoming data
func (r *instanceManager) startRuleChain(rulechainmodel *RuleChain) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.rulechains[rulechainmodel.ID]; found {
		logr.Debugf("rule chain '%s' is already started", rulechainmodel.ID)
		return nil
	}
	// create the internal runtime rulechain
	rulechain, errs := newRuleChainInstance(rulechainmodel.Channel, rulechainmodel.SubTopic, rulechainmodel.Payload)
	if len(errs) > 0 {
		return errs[0]
	}
	rulechainmodel.Status = RULE_STATUS_STARTED

	r.addInstanceInternal(rulechainmodel.ID, rulechain)
	return nil
}

// addInstanceInternal add a new rulechain instance internally with
// specified adaptor id
func (r *instanceManager) addInstanceInternal(rulechainID string, instance *ruleChainInstance) {
	r.rulechains[rulechainID] = instance
}

// stopRuleChain stop the rule chain
func (r *instanceManager) stopRuleChain(rulechainmodel *RuleChain) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.rulechains[rulechainmodel.ID]; !found {
		logr.Debugf("rule chain '%s' is not found", rulechainmodel.ID)
		return fmt.Errorf("rule chain '%s' no exist", rulechainmodel.ID)
	}
	delete(r.rulechains, rulechainmodel.ID)
	rulechainmodel.Status = RULE_STATUS_STOPPED
	return nil
}

// deleteRuleChain remove rule chain
func (c *instanceManager) deleteRuleChain(rulechain *RuleChain) error {
	return nil
}

func (c *instanceManager) HandleMessage(rulechainmessage message.Message, msg *mainflux.Message) error {
	for _, rulechaininstance := range c.rulechains {
		if rulechaininstance.channel == msg.GetChannel() && rulechaininstance.subTopic == msg.GetSubtopic() {
			if node, found := rulechaininstance.nodes[rulechaininstance.firstRuleNodeId]; found {
				go node.Handle(rulechainmessage)
			}
		}
	}
	return nil
}
