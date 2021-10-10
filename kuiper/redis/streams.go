package redis

import (
	"context"
	"fmt"

	"github.com/cloustone/pandas/kuiper"
	"github.com/go-redis/redis"
)

const (
//	keyPrefix = "stream_key"
//	idPrefix  = "stream"
)

var _ kuiper.StreamCache = (*streamCache)(nil)

type streamCache struct {
	client *redis.Client
}

// NewStreamCache returns redis stream cache implementation.
func NewStreamCache(client *redis.Client) kuiper.StreamCache {
	return &streamCache{
		client: client,
	}
}

func (tc *streamCache) Save(_ context.Context, streamKey string, streamID string) error {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, streamKey)
	if err := tc.client.Set(tkey, streamID, 0).Err(); err != nil {
		return err
	}

	tid := fmt.Sprintf("%s:%s", idPrefix, streamID)
	return tc.client.Set(tid, streamKey, 0).Err()
}

func (tc *streamCache) ID(_ context.Context, streamKey string) (string, error) {
	tkey := fmt.Sprintf("%s:%s", keyPrefix, streamKey)
	streamID, err := tc.client.Get(tkey).Result()
	if err != nil {
		return "", err
	}

	return streamID, nil
}

func (tc *streamCache) Remove(_ context.Context, streamID string) error {
	tid := fmt.Sprintf("%s:%s", idPrefix, streamID)
	key, err := tc.client.Get(tid).Result()
	if err != nil {
		return err
	}

	tkey := fmt.Sprintf("%s:%s", keyPrefix, key)

	return tc.client.Del(tkey, tid).Err()
}
