package redis

import (
	"context"
	"fmt"

	"github.com/cloustone/pandas/kuiper"
	"github.com/go-redis/redis"
)

const (
	keyPrefix = "rule_key"
	idPrefix  = "rule"
)

var _ kuiper.RuleCache = (*ruleCache)(nil)

type ruleCache struct {
	client *redis.Client
}

// NewRuleCache returns redis rule cache implementation.
func NewRuleCache(client *redis.Client) kuiper.RuleCache {
	return &ruleCache{
		client: client,
	}
}

func (tc *ruleCache) Save(_ context.Context, ruleKey string, ruleID string) error {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, ruleKey)
	if err := tc.client.Set(tkey, ruleID, 0).Err(); err != nil {
		return err
	}

	tid := fmt.Sprintf("%s:%s", idPrefix, ruleID)
	return tc.client.Set(tid, ruleKey, 0).Err()
}

func (tc *ruleCache) ID(_ context.Context, ruleKey string) (string, error) {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, ruleKey)
	ruleID, err := tc.client.Get(tkey).Result()
	if err != nil {
		return "", err
	}

	return ruleID, nil
}

func (tc *ruleCache) Remove(_ context.Context, ruleID string) error {
	tid := fmt.Sprintf("%s:%s", idPrefix, ruleID)
	key, err := tc.client.Get(tid).Result()
	if err != nil {
		return err
	}

	tkey := fmt.Sprintf("%s:%s", keyPrefix, key)

	return tc.client.Del(tkey, tid).Err()
}
