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
package users

import "time"

type Principal struct {
	ID             string
	Token          string
	Username       string
	Password       string
	Roles          []string `gorm:"type:string[]"`
	PhoneNumbers   string
	LastMFA        string
	LastMFAUpdated time.Time
}

func NewPrincipal(username, pwd string) *Principal {
	return &Principal{
		Roles: []string{},
	}
}

func (p *Principal) WithRole(role string) *Principal {
	p.Roles = append(p.Roles, role)
	return p
}

func (p *Principal) WithRoles(roles ...string) *Principal {
	p.Roles = append(p.Roles, roles...)
	return p
}

type PrincipalDefinition struct{}