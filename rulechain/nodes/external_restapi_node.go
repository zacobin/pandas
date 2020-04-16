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

type externalRestapiNode struct {
	bareNode
	RestEndpointUrlPattern         string
	RequestMethod                  string
	headers                        map[string]string
	UseSimpleClientHttpFactory     bool
	ReadTimeoutMs                  int
	MaxParallelRequestsCount       int
	UseRedisQueueForMsgPersistence bool
	trimQueue                      bool
	MaxQueueSize                   int
}

type externalRestapiNodeFactory struct{}

func (f externalRestapiNodeFactory) Name() string     { return "ExternalRestapiNode" }
func (f externalRestapiNodeFactory) Category() string { return NODE_CATEGORY_EXTERNAL }
func (f externalRestapiNodeFactory) Create(id string, meta Metadata) (Node, error) {
	labels := []string{"True", "False"}
	node := &externalRestapiNode{
		bareNode: newBareNode(f.Name(), id, meta, labels),
	}
	return decodePath(meta, node)
}
