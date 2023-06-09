package controller

import (
	"GoRedisLearn/RedisUtil"
	"testing"
)

func TestHyper(t *testing.T) {
	rds := RedisUtil.RedisUtil

	rds.SET("caojinbo", "lxh", 0)
}
