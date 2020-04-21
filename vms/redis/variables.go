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
	modelKeyPrefix = "model_key"
	modelIdPrefix  = "model"
)

var _ vms.ModelCache = (*modelCache)(nil)

type modelCache struct {
	client *redis.Client
}

// NewModelCache returns redis model cache implementation.
func NewModelCache(client *redis.Client) vms.ModelCache {
	return &modelCache{
		client: client,
	}
}

func (tc *modelCache) Save(_ context.Context, modelKey string, modelID string) error {
	tkey := fmt.Sprintf("%s:%s", modelKeyPrefix, modelKey)
	if err := tc.client.Set(tkey, modelID, 0).Err(); err != nil {
		return err
	}

	tid := fmt.Sprintf("%s:%s", modelIdPrefix, modelID)
	return tc.client.Set(tid, modelKey, 0).Err()
}

func (tc *modelCache) ID(_ context.Context, modelKey string) (string, error) {
	tkey := fmt.Sprintf("%s:%s", modelKeyPrefix, modelKey)
	modelID, err := tc.client.Get(tkey).Result()
	if err != nil {
		return "", err
	}

	return modelID, nil
}

func (tc *modelCache) Remove(_ context.Context, modelID string) error {
	tid := fmt.Sprintf("%s:%s", modelIdPrefix, modelID)
	key, err := tc.client.Get(tid).Result()
	if err != nil {
		return err
	}

	tkey := fmt.Sprintf("%s:%s", modelKeyPrefix, key)

	return tc.client.Del(tkey, tid).Err()
}
