package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"GoRedisLearn/response"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

var ctx = context.Background()

func QueryById(c *gin.Context) {
	id := c.PostForm("id")
	key := "cache:shop:" + id

	db := DB.GetDB()
	rds := RedisUtil.RedisUtil
	// 1.从redis查询商铺缓存
	shop, _ := rds.GET(key)
	var e model.TbShop

	// 2.判断是否存在
	if shop != "" && shop != "null" {
		// 3.存在直接返回
		err := json.Unmarshal([]byte(shop), &e)
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{"data": e})
		return
	}
	// 命中为空 返回错误信息
	if shop == "null" {
		response.Fail(c, gin.H{"msg": "用户信息不存在"})
	}
	// 4.不存在，根据id查询数据库
	tx := db.Where("id=?", id).Find(&e)
	if tx.RowsAffected == 0 {
		// 5.不存在，返回错误
		err := rds.SET(key, nil, 3*time.Minute)
		if err != nil {
			panic(err)
		}
		response.Fail(c, nil)
		return
	}
	// 6.存在 写入redis 设置过期时间
	err := rds.SET(key, e, 30*time.Minute)
	if err != nil {
		panic(err)
	}
	// 7.返回
	response.Success(c, gin.H{"data": e})
}

func Update(c *gin.Context) {
	var shop = model.TbShop{
		Id:   2,
		Name: "caojinbo",
	}
	err := c.ShouldBind(&shop)
	if err != nil {
		panic(err)
	}

	db := DB.GetDB()
	rds := RedisUtil.RedisUtil

	// update database
	db.Debug().Updates(shop)

	// delete cache
	err = rds.DEL("cache:shop:" + strconv.Itoa(shop.Id))
	if err != nil {
		panic(err)
	}

}
