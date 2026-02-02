package infra

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository struct {
	Client *redis.Client
}

func (r *RedisRepository) Set(ctx context.Context, key string, value string) error {
	return r.Client.Set(ctx, key, value, 10*time.Minute).Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
