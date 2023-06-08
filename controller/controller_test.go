package controller

import (
	"GoRedisLearn/RedisUtil"
	"fmt"
	"testing"
)

func TestHyper(t *testing.T) {
	rds := RedisUtil.RedisUtil

	for i := 0; i < 1000000; i++ {
		rds.PFAdd(RCTX, "hl2", i)
	}
	fmt.Println("values:", rds.PFCount(RCTX, "hl2").Val())
}
