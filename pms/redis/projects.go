// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"

	"github.com/cloustone/pandas/pms"
	"github.com/go-redis/redis"
)

const (
	keyPrefix = "project_key"
	idPrefix  = "project"
)

var _ pms.ProjectCache = (*projectCache)(nil)

type projectCache struct {
	client *redis.Client
}

// NewProjectCache returns redis project cache implementation.
func NewProjectCache(client *redis.Client) pms.ProjectCache {
	return &projectCache{
		client: client,
	}
}

func (tc *projectCache) Save(_ context.Context, projectKey string, projectID string) error {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, projectKey)
	if err := tc.client.Set(tkey, projectID, 0).Err(); err != nil {
		return err
	}

	tid := fmt.Sprintf("%s:%s", idPrefix, projectID)
	return tc.client.Set(tid, projectKey, 0).Err()
}

func (tc *projectCache) ID(_ context.Context, projectKey string) (string, error) {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, projectKey)
	projectID, err := tc.client.Get(tkey).Result()
	if err != nil {
		return "", err
	}

	return projectID, nil
}

func (tc *projectCache) Remove(_ context.Context, projectID string) error {
	tid := fmt.Sprintf("%s:%s", idPrefix, projectID)
	key, err := tc.client.Get(tid).Result()
	if err != nil {
		return err
	}

	tkey := fmt.Sprintf("%s:%s", keyPrefix, key)

	return tc.client.Del(tkey, tid).Err()
}
