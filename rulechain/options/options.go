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
package options

import (
	broadcast_options "github.com/cloustone/pandas/pkg/broadcast"
	genericoptions "github.com/cloustone/pandas/pkg/server/options"
	"github.com/spf13/pflag"
)

type ServingOptions struct {
	SecureServing    *genericoptions.SecureServingOptions
	BroadcastServing *broadcast_options.ServingOptions
	DeployMode       string
	EtcdConnectURL   string
}

func NewServingOptions() *ServingOptions {
	s := ServingOptions{
		SecureServing:    genericoptions.NewSecureServingOptions("dmms"),
		BroadcastServing: broadcast_options.NewServingOptions(),
		DeployMode:       "local",
	}
	return &s
}

func (s *ServingOptions) AddFlags(fs *pflag.FlagSet) {
	s.SecureServing.AddFlags(fs)
	s.BroadcastServing.AddFlags(fs)
	fs.StringVar(&s.DeployMode, "deploy-mode", s.DeployMode, "If empty, use local deploy mode")
	fs.StringVar(&s.EtcdConnectURL, "etcd-connect-url", s.EtcdConnectURL, "If non empty, deployed as rulechain cluster")
}

func (s *ServingOptions) IsNoneLocalMode() bool {
	return s.DeployMode == "etcd" && s.EtcdConnectURL != ""
}
