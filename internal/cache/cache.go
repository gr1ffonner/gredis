package cache

import (
	"context"
	"gredis/pkg/logging"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	logger logging.Logger
}

func NewRedicClient(logger logging.Logger) (*RedisClient, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "redispass",
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Redis connected")

	return &RedisClient{
		client: client,
		logger: logger,
	}, nil
}

func (r *RedisClient) Set(key string, value string) error {
	ctx := context.Background()

	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		r.logger.Fatal(err)
	}
	return nil
}

func (r *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		r.logger.Fatal(err)
	}
	return val, nil
}
