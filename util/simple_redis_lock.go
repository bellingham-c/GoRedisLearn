package util

import (
	"GoRedisLearn/RedisUtil"
	"context"
	"time"
)

var (
	KEY     = "lock:oversold"
	RDS_CTX = context.Background()
	rds     = RedisUtil.GetClient()
)

func TryLock(UUID string, c context.Context) bool {
	// 设置失效时间，保证锁失效机制，防止死锁
	t := 30 * time.Second
	nx := rds.SetNX(c, KEY, UUID, t)
	rds.Set(c, "caojinbo", UUID, 1*time.Minute)
	return nx.Val()
}

func UnLock(UUID string, c context.Context) bool {
	t_uuid := rds.Get(c, KEY).Val()
	if UUID == t_uuid {
		del := rds.Del(RDS_CTX, KEY).Val()
		if del > 0 {
			return true
		}
	}
	return false
}
