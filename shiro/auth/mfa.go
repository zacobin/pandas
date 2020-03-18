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
package auth

import (
	"errors"

	"github.com/cloustone/pandas/shiro/auth/sms"
	"github.com/cloustone/pandas/shiro/options"
	"github.com/cloustone/pandas/shiro/realms"
)

type MFAuthenticator interface {
	Authenticate(principal *realms.Principal, factor ...string) error
}

func NewMFAuthenticator(servingOptions *options.ServingOptions) MFAuthenticator {
	switch servingOptions.MFA {
	default:
		return &defaultMFA{
			smsAuthenticator: NewSmsAuthenticator(servingOptions),
		}
	}
}

type defaultMFA struct {
	smsAuthenticator sms.Authenticator
}

// TODO: how to add customized template params here in future
const (
	templateCode  = "SMS_82045083"
	templateParam = "{\"code\":\"1234\"}"
)

func (m *defaultMFA) Authenticate(principal *realms.Principal, factor ...string) error {
	if len(factor) == 0 {
		return errors.New("invalid arguments")
	}
	_, err := m.smsAuthenticator.Execute(principal.PhoneNumbers, factor[0], templateCode, templateParam)
	return err
}
