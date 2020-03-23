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

// The response code which stands for a sms is sent successfully.
const ResponseCodeOk = "OK"

// The Response of sending sms API.
type Response struct {
	RawResponse []byte `json:"-"` // The raw response from server.
	RequestId   string `json:"RequestId"`
	Code        string `json:"Code"`
	Message     string `json:"Message"`
	BizId       string `json:"BizId"`
}

func (m *Response) IsSuccessful() bool {
	return m.Code == ResponseCodeOk
}
