package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	TTL = 10 * time.Second
)

type RedisConfig struct {
	Network    string `yaml:"network" env:"NETWORK" env-default:"tcp"`
	Host       string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port       string `yaml:"port" env:"PORT" env-default:"6379"`
	Password   string `yaml:"password" env:"PASSWORD" env-default:""`
	Database   int    `yaml:"database" env:"DATABASE" env-default:"0"`
	MaxRetries int    `yaml:"max_retries" env:"MAX_RETRIES" env-default:"3"`
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg RedisConfig) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:   cfg.Password,
		DB:         cfg.Database,
		MaxRetries: cfg.MaxRetries,
	})

	return &RedisCache{
		client: rdb,
	}
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cmd := r.client.Set(ctx, key, value, expiration)

	if _, err := cmd.Result(); err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	cmd := r.client.Get(ctx, key)

	result, err := cmd.Result()
	if err != nil {
		return "", err
	}

	return result, nil
}
