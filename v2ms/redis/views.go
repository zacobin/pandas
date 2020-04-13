// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"

	"github.com/cloustone/pandas/v2ms"
	"github.com/go-redis/redis"
)

const (
	viewKeyPrefix = "view_key"
	viewIdPrefix  = "view"
)

var _ v2ms.ViewCache = (*viewCache)(nil)

type viewCache struct {
	client *redis.Client
}

// NewViewCache returns redis view cache implementation.
func NewViewCache(client *redis.Client) v2ms.ViewCache {
	return &viewCache{
		client: client,
	}
}

func (tc *viewCache) Save(_ context.Context, viewKey string, viewID string) error {
	tkey := fmt.Sprintf("%s:%s", viewKeyPrefix, viewKey)
	if err := tc.client.Set(tkey, viewID, 0).Err(); err != nil {
		return err
	}

	tid := fmt.Sprintf("%s:%s", viewIdPrefix, viewID)
	return tc.client.Set(tid, viewKey, 0).Err()
}

func (tc *viewCache) ID(_ context.Context, viewKey string) (string, error) {
	tkey := fmt.Sprintf("%s:%s", viewKeyPrefix, viewKey)
	viewID, err := tc.client.Get(tkey).Result()
	if err != nil {
		return "", err
	}

	return viewID, nil
}

func (tc *viewCache) Remove(_ context.Context, viewID string) error {
	tid := fmt.Sprintf("%s:%s", viewIdPrefix, viewID)
	key, err := tc.client.Get(tid).Result()
	if err != nil {
		return err
	}

	tkey := fmt.Sprintf("%s:%s", viewKeyPrefix, key)

	return tc.client.Del(tkey, tid).Err()
}
