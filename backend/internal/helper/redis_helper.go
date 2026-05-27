package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const DefaultCacheTTL = 10 * time.Minute

func CacheKey(pattern string, args ...any) string {
	return fmt.Sprintf(pattern, args...)
}

func GetOrSet[T any](ctx context.Context, client *redis.Client, key string, ttl time.Duration, fetch func() (T, error)) (T, error) {
	val, err := client.Get(ctx, key).Bytes()
	if err == nil {
		var result T
		if jsonErr := json.Unmarshal(val, &result); jsonErr == nil {
			return result, nil
		}
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		// Redis error — fall through to DB
	}

	result, err := fetch()
	if err != nil {
		return result, err
	}

	if data, jsonErr := json.Marshal(result); jsonErr == nil {
		client.Set(ctx, key, data, ttl)
	}

	return result, nil
}

func DeleteCache(ctx context.Context, client *redis.Client, keys ...string) {
	client.Del(ctx, keys...)
}
