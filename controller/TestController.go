package controller

import (
	"GoRedisLearn/RedisUtil"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Test(c *gin.Context) {
	//var ctx = context.Background()
	//rds := RedisUtil.GetClient()
	//rds.Set(ctx, "test", nil, 30*time.Minute)
	//
	//res := rds.Get(ctx, "test").Val()
	rds := RedisUtil.RedisUtil
	err := rds.SETWithOutJson("test", nil, 30*time.Minute)
	if err != nil {
		panic(err)
	}
	res, err := rds.GET("test")

	if res == "" {
		fmt.Println("res", res)
	}
	fmt.Println("res123", res)

}
