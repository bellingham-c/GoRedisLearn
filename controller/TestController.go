package controller

import (
	"GoRedisLearn/RedisUtil"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func Test(c *gin.Context) {
	var ctx = context.Background()
	rds := RedisUtil.GetClient()
	// Define Lua Script
	script, err := os.ReadFile("./lua/a.lua")
	fmt.Println("script", string(script))
	// Call EVAL command
	err = rds.Eval(ctx, string(script), []string{"capjinbo"}, "value").Err()
	if err != nil {
		// handle error
		//panic(err)
	} else {
		fmt.Println("SET success")
	}
	res := rds.Get(ctx, "caojinbo").Val()
	fmt.Println("res", res)
	// caojinceshi
}
