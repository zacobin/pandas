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
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap"
)

const (
	LdapAdaptorName = "ldap"
)

type ldapRealmProvider struct {
	conn  *ldap.Conn
	realm Realm
}

func newLdapRealmProvider(r Realm) (RealmProvider, error) {
	conn, err := ldap.Dial("tcp", r.ServiceConnectURL)
	if err != nil {
		return nil, err
	}

	err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}

	err = conn.Bind(r.Username, r.Password)
	if err != nil {
		return nil, err
	}

	return &ldapRealmProvider{conn: conn, realm: r}, nil
}

func (l *ldapRealmProvider) Authenticate(principal Principal) error {
	searchRequest := ldap.NewSearchRequest(
		l.realm.SearchDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=inetOrgPerson)(mail=%s))", principal.Username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return fmt.Errorf("User does not exist or too many entries returned")
	}

	userDN := sr.Entries[0].DN
	err = l.conn.Bind(userDN, principal.Password)
	if err != nil {
		return err
	}

	err = l.conn.Bind(l.realm.Username, l.realm.Password)
	if err != nil {
		return err
	}

	return nil
}
