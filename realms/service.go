// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package realms

import (
	"context"
	"sync"

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/errors"
)

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict = errors.New("email already taken")

	// ErrMalformedEntity indicates malformed entity specification
	// (e.g. invalid realmname or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrRealmNotFound indicates a non-existent realm request.
	ErrRealmNotFound = errors.New("non-existent realm")

	// ErrScanMetadata indicates problem with metadata in db.
	ErrScanMetadata = errors.New("Failed to scan metadata")

	// ErrMissingEmail indicates missing email for password reset request.
	ErrMissingEmail = errors.New("missing email for password reset")

	// ErrUnauthorizedPrincipal indicate the pricipal can not be recognized
	ErrUnauthorizedPrincipal = errors.New("unauthorized principal")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// Register creates new realm . In case of the failed registration, a
	// non-nil error value is returned.
	Register(ctx context.Context, realm Realm) error

	// Get authenticated realm info for the given token.
	RealmInfo(ctx context.Context, token, name string) (Realm, error)

	// UpdateRealm updates the realm metadata.
	UpdateRealm(ctx context.Context, token string, realm Realm) error

	// RevokeRealm remove a realm
	RevokeRealm(ctx context.Context, token, name string) error

	// ListRealm return all realms info
	ListRealms(ctx context.Context, token string) ([]Realm, error)

	// Identify identify whether the principal is legal
	Identify(ctx context.Context, token string, principal Principal) error
}

var _ Service = (*realmService)(nil)

type realmService struct {
	realms    RealmRepository
	hasher    Hasher
	auth      mainflux.AuthNServiceClient
	providers []RealmProvider
	mutex     sync.RWMutex
}

// New instantiates the realms service implementation
func New(realms RealmRepository, hasher Hasher, auth mainflux.AuthNServiceClient, providers []RealmProvider) Service {
	return &realmService{
		realms:    realms,
		hasher:    hasher,
		auth:      auth,
		providers: providers,
		mutex:     sync.RWMutex{},
	}
}

func (svc realmService) Register(ctx context.Context, realm Realm) error {
	hash, err := svc.hasher.Hash(realm.Password)
	if err != nil {
		return errors.Wrap(ErrMalformedEntity, err)
	}

	realm.Password = hash
	return svc.realms.Save(ctx, realm)
}

func (svc realmService) RealmInfo(ctx context.Context, token, name string) (Realm, error) {
	_, err := svc.identify(ctx, token)
	if err != nil {
		return Realm{}, err
	}

	realm, err := svc.realms.Retrieve(ctx, name)
	if err != nil {
		return Realm{}, errors.Wrap(ErrUnauthorizedAccess, err)
	}

	return realm, nil
}

func (svc realmService) UpdateRealm(ctx context.Context, token string, r Realm) error {
	_, err := svc.identify(ctx, token)
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.realms.Update(ctx, r)
}

func (svc realmService) ListRealms(ctx context.Context, token string) ([]Realm, error) {
	_, err := svc.identify(ctx, token)
	if err != nil {
		return nil, errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.realms.List(ctx)
}

func (svc realmService) RevokeRealm(ctx context.Context, token string, name string) error {
	_, err := svc.identify(ctx, token)
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.realms.Revoke(ctx, name)
}

func (svc realmService) Identify(ctx context.Context, token string, principal Principal) error {
	_, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
	}
	svc.mutex.Lock()
	defer svc.mutex.Unlock()
	for _, provider := range svc.providers {
		if err := provider.Authenticate(principal); err == nil {
			return nil
		}
	}
	return ErrUnauthorizedPrincipal
}

func (svc realmService) identify(ctx context.Context, token string) (string, error) {
	_, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return "", nil
}
