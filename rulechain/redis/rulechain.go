// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"

	"github.com/cloustone/pandas/rulechain"
	"github.com/go-redis/redis"
)

const (
	keyPrefix = "rulechain_key"
	idPrefix  = "rulechain"
)

var _ rulechain.RuleChainCache = (*rulechainCache)(nil)

type rulechainCache struct {
	client *redis.Client
}

// NewRuleChainCache returns redis rulechain cache implementation.
func NewRuleChainCache(client *redis.Client) rulechain.RuleChainCache {
	return &rulechainCache{
		client: client,
	}
}

func (tc *rulechainCache) Save(_ context.Context, rulechainKey string, rulechainID string) error {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, rulechainKey)
	if err := tc.client.Set(tkey, rulechainID, 0).Err(); err != nil {
		return err
	}

	tid := fmt.Sprintf("%s:%s", idPrefix, rulechainID)
	return tc.client.Set(tid, rulechainKey, 0).Err()
}

func (tc *rulechainCache) ID(_ context.Context, rulechainKey string) (string, error) {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, rulechainKey)
	rulechainID, err := tc.client.Get(tkey).Result()
	if err != nil {
		return "", err
	}

	return rulechainID, nil
}

func (tc *rulechainCache) Remove(_ context.Context, rulechainID string) error {
	tid := fmt.Sprintf("%s:%s", idPrefix, rulechainID)
	key, err := tc.client.Get(tid).Result()
	if err != nil {
		return err
	}

	tkey := fmt.Sprintf("%s:%s", keyPrefix, key)

	return tc.client.Del(tkey, tid).Err()
}
