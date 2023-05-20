package controller

import (
	"GoRedisLearn/RedisUtil"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func Test(c *gin.Context) {
	var ctx = context.Background()
	rds := RedisUtil.GetClient()
	rds.Set(ctx, "caojinbo", "test", 30*time.Second)
}
