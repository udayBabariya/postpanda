package database

import (
	"context"
	"fmt"
	"postpanda/backend-go/internal/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectRedis() (*redis.Client, error) {
	opts, err := redis.ParseURL(config.App.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	Redis = client
	return client, nil
}
