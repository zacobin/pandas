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
package authz

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// roleManifest contains builtin role and routes, the file can be loaded since
// the security manager is started, security manager use default role policy to
// initialize it
type roleManifest struct {
	Version string  `json:"version"`
	Roles   []*Role `json:"roles"`
}

// loadRoles load roles from manifest file
func loadRoles(fileName string) ([]*Role, error) {
	manifest := &roleManifest{
		Roles: []*Role{},
	}

	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(buf, manifest); err != nil {
		return nil, err
	}

	roleMaps := make(map[string]*Role)
	// check mainifest's validity
	for _, role := range manifest.Roles {
		if _, found := roleMaps[role.Name]; found {
			return nil, errors.New("invalid role manifest file")
		}
		roleMaps[role.Name] = role
	}
	return manifest.Roles, nil
}
