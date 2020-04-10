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

	"github.com/cloustone/pandas/apimachinery/models"
)

type SaveAttributesNode struct {
	bareNode
}

type saveAttributesNodeFactory struct{}

func (f saveAttributesNodeFactory) Name() string     { return "SaveAttributesNode" }
func (f saveAttributesNodeFactory) Category() string { return NODE_CATEGORY_ACTION }
func (f saveAttributesNodeFactory) Create(id string, meta Metadata) (Node, error) {
	labels := []string{"Success", "Failure"}
	node := &SaveAttributesNode{
		bareNode: newBareNode(f.Name(), id, meta, labels),
	}
	return decodePath(meta, node)
}

func (n *SaveAttributesNode) Handle(msg models.Message) error {
	successLableNode := n.GetLinkedNode("Success")
	failureLableNode := n.GetLinkedNode("Failure")
	if successLableNode == nil || failureLableNode == nil {
		return fmt.Errorf("no valid label linked node in %s", n.Name())
	}
	if msg.GetType() != "POST_ATTRIBUTES_REQUEST" {
		return failureLableNode.Handle(msg)
	}

	return nil
}
