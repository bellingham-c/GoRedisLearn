package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	db := DB.GetDB()
	rds := RedisUtil.RedisUtil

	var shop []model.TbShop
	db.Where("id<3").Find(&shop)

	fmt.Println("shop", shop)
	err := rds.HSET("test", shop)
	if err != nil {
		panic(err)
	}

	resultMap, err := rds.HGETALL("test")
	fmt.Println(resultMap)
}
