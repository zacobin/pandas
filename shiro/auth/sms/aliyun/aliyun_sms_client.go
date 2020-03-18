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
package aliyun

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cloustone/pandas/shiro/auth/sms"
	"github.com/cloustone/pandas/shiro/options"
)

type smsAuthenticator struct {
	request         *aliyunCommunicationRequest
	gatewayUrl      string
	accessKeyID     string
	accessKeySecret string
	client          *http.Client
}

func NewAuthenticator(servingOptions *options.ServingOptions) sms.Authenticator {
	return &smsAuthenticator{
		request:         &aliyunCommunicationRequest{},
		gatewayUrl:      servingOptions.SmsAccessURL,
		client:          &http.Client{},
		accessKeyID:     servingOptions.SmsAccessKeyID,
		accessKeySecret: servingOptions.SmsAccessKeySecret,
	}
}

func (s *smsAuthenticator) Execute(phoneNumbers, signName, templateCode, templateParam string) (*sms.Response, error) {
	err := s.request.SetParamsValue(s.accessKeyID, phoneNumbers, signName, templateCode, templateParam)
	if err != nil {
		return nil, err
	}
	endpoint, err := s.request.BuildSmsRequestEndpoint(s.accessKeySecret, s.gatewayUrl)
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("GET", endpoint, nil)
	response, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	result := new(sms.Response)
	err = json.Unmarshal(body, result)

	result.RawResponse = body
	return result, err
}
