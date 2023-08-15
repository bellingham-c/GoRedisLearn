package newRedisLock

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// Client Redis 客户端.
type Client struct {
	pool *redis.Client
}

func NewClient() *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "43.143.241.157:6379",
		Password: "123456",
		DB:       0,
	})

	return &Client{pool: client}
}

func (c *Client) SetNEX(ctx context.Context, key, value string, expireSeconds int64) (bool, error) {
	if key == "" || value == "" {
		return false, errors.New("redis SET keyNX or value can't be empty")
	}

	setNX := c.pool.SetNX(ctx, key, value, time.Duration(expireSeconds)*time.Second)
	return setNX.Result()
}

// Eval 支持使用 lua 脚本.
func (c *Client) Eval(ctx context.Context, script, key, token string, expireSeconds int64) (bool, error) {
	keys := make([]string, 1)
	argv := make([]interface{}, 2)

	keys[0] = key
	argv[0] = token
	argv[1] = expireSeconds

	eval := c.pool.Eval(ctx, script, keys, argv)
	if eval.Err() != nil {
		return false, eval.Err()
	}
	return true, nil
}
