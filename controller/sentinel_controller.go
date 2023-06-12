package controller

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type userCredentialLoginCount struct {
	ID                string // 用户id也可能是用户名
	account_pwd_count int    // 账号或密码错误次数
	token_count       int    // 令牌错误次数
}

func SentinelController() {
	ctx := context.Background()
	rdb := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			// shardName => host:port
			"shard1": "192.168.80.129:7000",
			"shard2": "192.168.80.129:7001",
			"shard3": "192.168.80.129:7002",
			"shard4": "192.168.80.129:8001",
			"shard5": "192.168.80.129:8002",
			"shard6": "192.168.80.129:8003",
		},
	})

	set := rdb.Set(ctx, "cao", "andong1", 0)
	fmt.Println("set", set)
}
