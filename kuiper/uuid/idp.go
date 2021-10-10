package uuid

import (
	"github.com/cloustone/pandas/kuiper"
	"github.com/gofrs/uuid"
)

var _ kuiper.IdentityProvider = (*uuidIdentityProvider)(nil)

type uuidIdentityProvider struct{}

// New instantiates a UUID identity provider.
func New() kuiper.IdentityProvider {
	return &uuidIdentityProvider{}
}

func (idp *uuidIdentityProvider) ID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
