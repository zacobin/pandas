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
package realms

import (
	"context"

	"github.com/sirupsen/logrus"
)

// RealmProvider represent entity that provide user identity authentication using ldap
type RealmProvider interface {
	Authenticate(principal Principal) error
}

// Realm
type Realm struct {
	Name              string `json:"name"`
	CertFile          string `json:"certFile"`
	KeyFile           string `json:"keyFile"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	ServiceConnectURL string `json:"serviceConnectURL"`
	SearchDN          string `json:"searchDN"`
}

// Validate returns an error if representtation is invalid
func (r Realm) Validate() error {
	if r.SearchDN == "" || r.ServiceConnectURL == "" {
		return ErrMalformedEntity
	}
	return nil
}

// RealmRepository specifies realm persistence API
type RealmRepository interface {
	// Save persists the realm
	Save(context.Context, Realm) error

	// Update updates the realm metdata
	Update(context.Context, Realm) error

	// Retrieve return realm by its identifier (i.e name)
	Retrieve(context.Context, string) (Realm, error)

	// RevokeRealm remove realm
	Revoke(context.Context, string) error

	// List return all reams
	List(context.Context) ([]Realm, error)
}

// NewRealProvider create a realm from specific options
func NewRealmProvider(realm Realm) (RealmProvider, error) {
	switch realm.Name {
	case LdapAdaptorName:
		realm, err := newLdapRealmProvider(realm)
		if err != nil {
			logrus.WithError(err).Errorf("invalid realm '%s' options", LdapAdaptorName)
			return nil, err
		}
		return realm, nil
	}
	logrus.Fatalf("invalid realm names")
	return nil, nil
}
