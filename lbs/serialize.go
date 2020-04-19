//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package lbs

import (
	"encoding/json"
	"errors"
)

var (
	JSONSerialization = SerializeOption{Format: "json"}
)

type SerializeOption struct {
	Format string
}

func Serialize(msg interface{}, opt SerializeOption) ([]byte, error) {
	switch opt {
	case JSONSerialization:
		return json.Marshal(msg)
	default:
		return nil, errors.New("invalid message codec")
	}
}

func Deserialize(buf []byte, opt SerializeOption, obj interface{}) error {
	switch opt {
	case JSONSerialization:
		return json.Unmarshal(buf, obj)
	default:
		return errors.New("invalid message codec")
	}
}