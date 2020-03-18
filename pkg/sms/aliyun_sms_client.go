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
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type aliyunSMS struct {
	request         *aliyunCommunicationRequest
	gatewayUrl      string
	accessKeyID     string
	accessKeySecret string
	client          *http.Client
}

func newAliyunSMS(servingOptions *ServingOptions) Client {
	return &aliyunSMS{
		request:         &aliyunCommunicationRequest{},
		gatewayUrl:      servingOptions.AccessURL,
		client:          &http.Client{},
		accessKeyID:     servingOptions.AccessKeyID,
		accessKeySecret: servingOptions.AccessKeySecret,
	}
}

func (s *aliyunSMS) Execute(phoneNumbers, signName, templateCode, templateParam string) (*Response, error) {
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

	result := new(Response)
	err = json.Unmarshal(body, result)

	result.RawResponse = body
	return result, err
}
