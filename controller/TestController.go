package controller

import (
	"GoRedisLearn/RedisUtil"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	token := c.GetHeader("Authorization")
	token = token[7:]
	rds := RedisUtil.RedisUtil
	res, err := rds.HGETALL(token)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
