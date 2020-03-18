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
	"errors"
	"time"

	"github.com/cloustone/pandas/pkg/sms"
	"github.com/cloustone/pandas/shiro/options"
	"github.com/cloustone/pandas/shiro/realms"
	"github.com/rs/xid"
)

const (
	MFADuration = 45
)

type MFAuthenticator interface {
	Notify(principal *realms.Principal) error
	Authenticate(principal *realms.Principal) error
}

func NewMFAuthenticator(servingOptions *options.ServingOptions, backstoreManager *backstoreManager) MFAuthenticator {
	switch servingOptions.MFA {
	default:
		return newDefaultMFA(servingOptions, backstoreManager)
	}
}

type defaultMFA struct {
	smsClient        sms.Client
	backstoreManager *backstoreManager
	servingOptions   *options.ServingOptions
}

func newDefaultMFA(servingOptions *options.ServingOptions, backstoreManager *backstoreManager) MFAuthenticator {
	return &defaultMFA{
		servingOptions:   servingOptions,
		smsClient:        sms.NewClient(servingOptions.SmsOptions),
		backstoreManager: backstoreManager,
	}
}

func (mfa *defaultMFA) Notify(principal *realms.Principal) error {
	signName := xid.New().String() // TODO: we should provide a readable text
	smsOptions := mfa.servingOptions.SmsOptions
	if _, err := mfa.smsClient.Execute(principal.PhoneNumbers, signName, smsOptions.TemplateCode, smsOptions.TemplateParam); err != nil {
		return err
	}
	principal.LastMFA = signName
	principal.LastMFAUpdated = time.Now()
	mfa.backstoreManager.updatePrincipal(principal)
	return nil
}

func (mfa *defaultMFA) Authenticate(principal *realms.Principal) error {
	lastMFA := principal.LastMFA
	if err := mfa.backstoreManager.getPrincipal(principal); err != nil {
		return err
	}
	duration := time.Since(principal.LastMFAUpdated)
	if lastMFA == principal.LastMFA && duration < MFADuration {
		return nil
	}
	return errors.New("invalid principal")
}
