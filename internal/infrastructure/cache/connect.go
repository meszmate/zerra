package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	*redis.Client
}

func New(redisURL string) (*Cache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Cache{
		Client: rdb,
	}, nil
}
