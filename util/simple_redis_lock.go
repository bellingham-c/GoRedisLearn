package util

import (
	"GoRedisLearn/RedisUtil"
	"context"
	"time"
)

var (
	KEY_PREFIX = "lock:"
	RDS_CTX    = context.Background()
	rds        = RedisUtil.GetClient()
)

func TryLock(UUID string, ttl int) bool {
	key := KEY_PREFIX + UUID
	t := time.Duration(ttl) * time.Minute
	nx := rds.SetNX(RDS_CTX, key, 1, t).Val()
	return nx
}

func UnLock(UUID string) bool {
	key := KEY_PREFIX + UUID
	del := rds.Del(RDS_CTX, key).Val()
	if del > 0 {
		return true
	}
	return false
}
