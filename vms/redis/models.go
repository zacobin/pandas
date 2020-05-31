// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"

	"github.com/cloustone/pandas/vms"
	"github.com/go-redis/redis"
)

const (
	keyPrefix = "variable_key"
	idPrefix  = "variable"
)

var _ vms.VariableCache = (*variableCache)(nil)

type variableCache struct {
	client *redis.Client
}

// NewVariableCache returns redis variable cache implementation.
func NewVariableCache(client *redis.Client) vms.VariableCache {
	return &variableCache{
		client: client,
	}
}

func (tc *variableCache) Save(_ context.Context, variableKey string, variableID string) error {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, variableKey)
	if err := tc.client.Set(tkey, variableID, 0).Err(); err != nil {
		return err
	}

	tid := fmt.Sprintf("%s:%s", idPrefix, variableID)
	return tc.client.Set(tid, variableKey, 0).Err()
}

func (tc *variableCache) ID(_ context.Context, variableKey string) (string, error) {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, variableKey)
	variableID, err := tc.client.Get(tkey).Result()
	if err != nil {
		return "", err
	}

	return variableID, nil
}

func (tc *variableCache) Remove(_ context.Context, variableID string) error {
	tid := fmt.Sprintf("%s:%s", idPrefix, variableID)
	key, err := tc.client.Get(tid).Result()
	if err != nil {
		return err
	}

	tkey := fmt.Sprintf("%s:%s", keyPrefix, key)

	return tc.client.Del(tkey, tid).Err()
}
