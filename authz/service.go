// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package authz

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
	// RetrieveRole return specified role
	RetrieveRole(ctx context.Context, token, roleName string) (Role, error)

	// RetrieveRoles return all default roles
	ListRoles(ctx context.Context, token string) ([]Role, error)

	// UpdateRole update specified role to change permissions
	UpdateRole(ctx context.Context, token string, role Role) error

	// Authorize checkwethere a role can access subject
	Authorize(ctx context.Context, token string, roleName string, subject Subject) error
}

var _ Service = (*authzService)(nil)

type authzService struct {
	repo   RoleRepository
	hasher Hasher
	auth   mainflux.AuthNServiceClient
	mutex  sync.RWMutex
}

// New instantiates the authz service implementation
func New(repo RoleRepository, hasher Hasher, auth mainflux.AuthNServiceClient) Service {
	return &authzService{
		repo:   repo,
		hasher: hasher,
		auth:   auth,
		mutex:  sync.RWMutex{},
	}
}

func (svc authzService) identify(ctx context.Context, token string) (string, error) {
	_, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return "", nil
}

// RetrieveRole return specified role
func (svc authzService) RetrieveRole(ctx context.Context, token, roleName string) (Role, error) {
	return Role{}, nil
}

// RetrieveRoles return all default roles
func (svc authzService) ListRoles(ctx context.Context, token string) ([]Role, error) {
	return []Role{}, nil
}

// UpdateRole update specified role to change permissions
func (svc authzService) UpdateRole(ctx context.Context, token string, role Role) error {
	return nil
}

// Authorize checkwethere a role can access subject
func (svc authzService) Authorize(ctx context.Context, token string, roleName string, subject Subject) error {
	return nil
}
