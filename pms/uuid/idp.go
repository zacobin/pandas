// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package uuid provides a UUID identity provider.
package uuid

import (
	"github.com/cloustone/pandas/v2ms"
	"github.com/gofrs/uuid"
)

var _ v2ms.IdentityProvider = (*uuidIdentityProvider)(nil)

type uuidIdentityProvider struct{}

// New instantiates a UUID identity provider.
func New() v2ms.IdentityProvider {
	return &uuidIdentityProvider{}
}

func (idp *uuidIdentityProvider) ID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (idp *uuidIdentityProvider) IsValid(u4 string) error {
	if _, err := uuid.FromString(u4); err != nil {
		return v2ms.ErrMalformedEntity
	}

	return nil
}
