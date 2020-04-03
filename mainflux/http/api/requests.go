// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import "github.com/cloustone/pandas/mainflux/broker"

type publishReq struct {
	msg   broker.Message
	token string
}
