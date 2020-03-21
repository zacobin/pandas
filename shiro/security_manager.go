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
	"sync"
	"time"

	"github.com/cloustone/pandas/shiro/auth"
	"github.com/cloustone/pandas/shiro/options"
	"github.com/cloustone/pandas/shiro/realms"
	. "github.com/cloustone/pandas/shiro/realms"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// SecurityManager is responsible for authenticate and simple authorization
type SecurityManager interface {
	UseAdaptor(Adaptor)
	AddDomainRealm(realms.Realm)
	Authenticate(principal *Principal, factor ...string) error
	Authorize(principal Principal, object *Object, action string) error
	GetAuthzDefinitions(principal Principal) ([]*AuthzDefinition, error)
	GetPrincipalDefinition(principal Principal) (*PrincipalDefinition, error)
	GetAllRoles() []*Role
	UpdateRole(r *Role) error
	UpdatePrincipal(principal Principal) error
}

// NewSecurityManager create security manager to hold all realms for
// authenticate
func NewSecurityManager(servingOptions *options.ServingOptions, backstoreManager *backstoreManager, mfa MFAuthenticator) SecurityManager {
	return newDefaultSecurityManager(servingOptions, backstoreManager, mfa)
}

// defaultSecuriityManager
type defaultSecurityManager struct {
	mutex            sync.RWMutex
	servingOptions   *options.ServingOptions
	backstoreManager *backstoreManager
	realms           []realms.Realm
	mfa              MFAuthenticator
}

// newDefaultSecurityManager return security manager instance
// All realms are created here, if failed, shiro must be restarted
func newDefaultSecurityManager(servingOptions *options.ServingOptions, backstoreManager *backstoreManager, mfa MFAuthenticator) *defaultSecurityManager {
	realms, err := NewRealmsWithFile(servingOptions.RealmConfigFile)
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	backstoreManager.loadRolesWithFile(servingOptions.RolesFile)
	return &defaultSecurityManager{
		mutex:            sync.RWMutex{},
		servingOptions:   servingOptions,
		backstoreManager: backstoreManager,
		realms:           realms,
		mfa:              mfa,
	}
}

// UseAdaptor use synchronization adaptor between shiro nodes
func (s *defaultSecurityManager) UseAdaptor(Adaptor) {}

// AddDomainRealm adds domain's specific realm
//realm is only a kind of interface you can initliaze it with ldaprealm so it will be a ldaprealm
func (s *defaultSecurityManager) AddDomainRealm(realm realms.Realm) {
	// TODO: add realm simply
	s.mutex.Lock()
	s.realms = append(s.realms, realm)
	s.mutex.Unlock()
}

// Authenticate iterate all realm to authenticate the principal
func (s *defaultSecurityManager) Authenticate(principal *Principal, factor ...string) error {
	authenticated := false

	for _, realm := range s.realms {
		if err := realm.Authenticate(principal); err == nil {
			if err := s.backstoreManager.getPrincipal(principal); err != nil {
				return errors.New("invalid user")
			}
			authenticated = true
			break
		}
	}
	if !authenticated {
		return errors.New("no valid realms")
	}

	// Two factor authentication
	if err := s.mfa.Authenticate(principal); err != nil {
		return err
	}

	claims := auth.JwtClaims{
		AccessId: principal.ID,
		Name:     principal.Username,
		Roles:    principal.Roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(6000)).Unix(),
			Issuer:    "pandas",
		},
	}
	token, err := auth.GenToken(claims)
	if err != nil {
		return err
	}
	principal.Token = token
	return nil
}

func (s *defaultSecurityManager) Authorize(principal Principal, object *Object, action string) error {
	return nil
}

// GetAllRoles return all builtin role's definitions
func (s *defaultSecurityManager) GetAllRoles() []*Role {
	return s.backstoreManager.getAllRoles()
}

// UpdateRole update a role's definition
func (s *defaultSecurityManager) UpdateRole(r *Role) error {
	return s.backstoreManager.updateRole(r)
}

func (s *defaultSecurityManager) GetAuthzDefinitions(principal Principal) ([]*AuthzDefinition, error) {
	return nil, nil
}
func (s *defaultSecurityManager) GetPrincipalDefinition(principal Principal) (*PrincipalDefinition, error) {
	return nil, nil
}

func (s *defaultSecurityManager) UpdatePrincipal(principal Principal) error {
	return s.backstoreManager.updatePrincipal(&principal)
}
