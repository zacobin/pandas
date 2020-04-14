// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"github.com/cloustone/pandas/lbs"
)

type listCollectionsReq struct {
	token string
}

func (req listCollectionsReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	return nil
}
