package util

import (
	"GoRedisLearn/RedisUtil"
	"context"
)

var (
	KEY_PREFIX = "lock:"
	ctx        = context.Background()
	rds        = RedisUtil.RedisUtil
)

func TryLock(name string, ttl int) {
	key := KEY_PREFIX + name

}

func UnLock() {

}
