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

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cloustone/pandas/users/realms"
	"github.com/cloustone/pandas/users/realms/ldap"
	"github.com/sirupsen/logrus"
)

type realmoptions struct {
	Realmoptions []realms.RealmOptions `json: realmoptions`
}

// NewReal create a realm from specific options
func NewRealm(realmOptions *realms.RealmOptions) (realms.Realm, error) {
	switch realmOptions.Name {
	case ldap.AdaptorName:
		realm, err := ldap.NewLdapRealm(realmOptions)
		if err != nil {
			logrus.WithError(err).Errorf("invalid realm '%s' options", ldap.AdaptorName)
			return nil, err
		}
		return realm, nil
	}
	logrus.Fatalf("invalid realm names")
	return nil, nil
}

// NewRealmsWithFile create realms from realms config file
func NewRealmsWithFile(fullFilePath string) ([]realms.Realm, error) {
	buf, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		logrus.WithError(err).Fatalf("open realms config file failed")
		return nil, err
	}
	var realmop realmoptions
	if err := json.Unmarshal(buf, &realmop); err != nil {
		logrus.WithError(err).Fatalf("illegal realm config file")
		return nil, err
	}

	realms := []realms.Realm{}
	for _, option := range realmop.Realmoptions {
		if realm, err := NewRealm(&option); err != nil {
			logrus.WithError(err)
			continue
		} else {
			realms = append(realms, realm)
		}
	}
	return realms, nil
}
