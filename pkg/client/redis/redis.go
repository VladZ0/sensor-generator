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
	Network         string `yaml:"network" env:"NETWORK" env-default:"tcp"`
	Host            string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port            string `yaml:"port" env:"PORT" env-default:"6379"`
	Password        string `yaml:"password" env:"PASSWORD" env-default:""`
	Database        int    `yaml:"database" env:"DATABASE" env-default:"0"`
	MaxRetries      int    `yaml:"max_retries" env:"MAX_RETRIES" env-default:"3"`
	DialTimeout     int    `yaml:"dial_timeout" env:"DIAL_TIMEOUT" env-default:"5"`
	ReadTimeout     int    `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"3"`
	WriteTimeout    int    `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"3"`
	PoolSize        int    `yaml:"pool_size" env:"POOL_SIZE" env-default:"10"`
	PoolTimeout     int    `yaml:"pool_timeout" env:"POOL_TIMEOUT" env-default:"0"`
	MinIdleConns    int    `yaml:"min_idle_conns" env:"MIN_IDLE_CONNS" env-default:"0"`
	MaxIdleConns    int    `yaml:"max_idle_conns" env:"MAX_IDLE_CONNS" env-default:"0"`
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time" env:"CONN_MAX_IDLE_TIME" env-default:"1800"`
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg RedisConfig) *RedisCache {
	if cfg.PoolTimeout == 0 {
		cfg.PoolTimeout = cfg.ReadTimeout + 1
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.Database,
		MaxRetries:      cfg.MaxRetries,
		DialTimeout:     time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(cfg.WriteTimeout) * time.Second,
		PoolSize:        cfg.PoolSize,
		PoolTimeout:     time.Duration(cfg.PoolTimeout) * time.Second,
		MinIdleConns:    cfg.MinIdleConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxIdleTime: time.Duration(cfg.ConnMaxIdleTime) * time.Second,
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
