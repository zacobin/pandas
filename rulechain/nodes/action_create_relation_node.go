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
	"github.com/cloustone/pandas/apimachinery/models"
)

type createRelationNode struct {
	bareNode
	Direction                       string
	RelationType                    string
	EntityType                      string
	EntityNamePattern               string
	EntityTypePattern               string
	EntityCacheExpiration           int64
	CreateEntityIfNotExists         bool
	ChangeOriginatorToRelatedEntity bool
	RemoveCurrentRelations          bool
}

type createRelationNodeFactory struct{}

func (f createRelationNodeFactory) Name() string     { return "CreateRelationNode" }
func (f createRelationNodeFactory) Category() string { return NODE_CATEGORY_ACTION }
func (f createRelationNodeFactory) Create(id string, meta Metadata) (Node, error) {
	labels := []string{"Success", "Failure"}
	node := &createRelationNode{
		bareNode: newBareNode(f.Name(), id, meta, labels),
	}
	return decodePath(meta, node)
}

func (n *createRelationNode) Handle(msg models.Message) error {
	return nil
}
