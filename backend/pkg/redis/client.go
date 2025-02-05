package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client wraps a Redis client instance.
type Client struct {
	Conn *redis.Client
	Ctx  context.Context
}

// NewClient creates and returns a new Redis client.
func NewClient(addr, password string, db int, timeout time.Duration) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("failed to connect to Redis: " + err.Error())
	}

	return &Client{
		Conn: rdb,
		Ctx:  context.Background(),
	}
}
