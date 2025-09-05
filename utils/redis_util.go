package utils

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

// 通用存储函数
func SetRedisStruct[T any](ctx context.Context, rdb redis.UniversalClient, key string, value T, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, jsonData, expiration).Err()
}

// 通用获取函数
func GetRedisStruct[T any](ctx context.Context, rdb redis.UniversalClient, key string) (T, error) {
	var zero T
	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return zero, err
	}

	var result T
	err = json.Unmarshal(data, &result)
	if err != nil {
		return zero, err
	}

	return result, nil
}
