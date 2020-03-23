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
package sms

import (
	"github.com/spf13/pflag"
)

const (
	templateCode  = "SMS_82045083"
	templateParam = "{\"code\":\"1234\"}"
)

type ServingOptions struct {
	AccessURL       string
	AccessKeyID     string
	AccessKeySecret string
	TemplateCode    string
	TemplateParam   string
}

func NewServingOptions() *ServingOptions {
	return &ServingOptions{
		TemplateCode:  templateCode,
		TemplateParam: templateParam,
	}
}

func (s *ServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.AccessURL, "sms-access-url", s.AccessURL, "sms authenticator access key url")
	fs.StringVar(&s.AccessKeyID, "sms-access-key-id", s.AccessKeyID, "sms authenticator access key id")
	fs.StringVar(&s.AccessKeySecret, "sms-access-key-secret", s.AccessKeySecret, "sms authenticator access key secret")
	fs.StringVar(&s.TemplateCode, "sms-template-code", s.TemplateCode, "sms template code")
	fs.StringVar(&s.TemplateParam, "sms-template-param", s.TemplateParam, "sms template param")
}
