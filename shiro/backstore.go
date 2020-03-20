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
package shiro

import (
	"github.com/cloustone/pandas/shiro/options"
	"github.com/cloustone/pandas/shiro/realms"
	"github.com/sirupsen/logrus"
)

type backstoreManager struct{}

func newBackstoreManager(servingOptions *options.ServingOptions) *backstoreManager {
	return &backstoreManager{}
}

func (b *backstoreManager) getPrincipal(pricipal *realms.Principal) error {
	// TODO: access database to retrieve principal's roles
	return nil
}

func (b *backstoreManager) updatePrincipal(principal *realms.Principal) error {
	return nil
}

func (b *backstoreManager) loadRoles(rolesFile string) {
	// Load builtin role's definitions
	_, err := loadRoles(rolesFile)
	if err != nil {
		logrus.WithError(err)
	}
}

func (b *backstoreManager) getAllRoles() []*Role {
	return nil
}

func (b *backstoreManager) updateRole(r *Role) {
}
