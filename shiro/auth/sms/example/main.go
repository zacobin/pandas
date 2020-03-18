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

package main

import (
	"encoding/json"
	"fmt"

	"github.com/cloustone/pandas/shiro/options"
)

var (
	gatewayUrl      = "http://dysmsapi.aliyuncs.com/"
	accessKeyId     = "LTAIbTnPbawglLIQ"
	accessKeySecret = ""
	phoneNumbers    = "13544285**2"
	signName        = "jenson"
	templateCode    = "SMS_82045083"
	templateParam   = "{\"code\":\"1234\"}"
)

func main() {
	sevingOptons := &options.ServingOptions{
		SmsAccessURL:       gatewayUrl,
		SmsAccessKeyID:     accessKeyID,
		SmsAccessKeySecret: accessKeySecret,
	}
	smsAuthc := aliyun.NewSmsAuthenticator(servingOptions)
	result, err := smsAuthc.Execute(phoneNumbers, signName, templateCode, templateParam)
	fmt.Println("Got raw response from server:", string(result.RawResponse))
	if err != nil {
		panic("Failed to send Message: " + err.Error())
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	if result.IsSuccessful() {
		fmt.Println("A SMS is sent successfully:", resultJson)
	} else {
		fmt.Println("Failed to send a SMS:", resultJson)
	}
}
