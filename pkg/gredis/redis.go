package gredis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func SetUp() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:            setting.RedisSetting.Addr,
		Password:        setting.RedisSetting.Password,
		PoolSize:        setting.RedisSetting.MaxActive,
		MinIdleConns:    setting.RedisSetting.MaxIdle,
		ConnMaxIdleTime: setting.RedisSetting.IdleTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func Set(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = RedisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func Exists(ctx context.Context, key string) (bool, error) {
	exists, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func Get(ctx context.Context, key string) ([]byte, error) {
	value, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return value, err
}

func Delete(ctx context.Context, key string) error {
	err := RedisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func LikeDeletes(ctx context.Context, pattern string) error {
	keys, err := RedisClient.Keys(ctx, "*"+pattern+"*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		err = RedisClient.Del(ctx, keys...).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
