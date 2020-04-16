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

import "github.com/cloustone/pandas/rulechain/message"

type ScriptEngine interface {
	ScriptOnMessage(msg message.Message, script string) (message.Message, error)
	//used by filter_switch_node
	ScriptOnSwitch(msg message.Message, script string) ([]string, error)
	//used by filter_script_node
	ScriptOnFilter(msg message.Message, script string) (bool, error)
	ScriptToString(msg message.Message, script string) (string, error)
}

func NewScriptEngine() ScriptEngine {
	return nil
}
