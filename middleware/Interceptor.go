package middleware

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"time"
)

var user model.TbUser

// RefreshTokenMiddleware 拦截所有请求，对token进行刷新，也就是对redis里面的token时效进行重置，三十分钟
func RefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1 获取请求头的token
		token := c.GetHeader("Authorization") //
		if token == "" {
			//token 不存在 放行
			return
		}
		//2 基于token获取redis中的用户信息
		token = token[7:]
		rds := RedisUtil.RedisUtil
		res, err := rds.HGETALL(token)
		if err != nil {
			panic(err)
		}
		//3 判断用户是否有效
		db := DB.GetDB()
		tx := db.Where("id=?", res["Id"]).Find(&user)
		if tx.RowsAffected == 0 {
			//用户不存在 放行
			return
		}
		//4 用户存在 刷新token有效期
		err = rds.EXPIRE(token, 30*time.Minute)
		if err != nil {
			panic(err)
		}
		// 放行
		// 什么都不用做 直接执行下一步
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("middleware auth")
		// 没有用户信息 拦截
		if reflect.DeepEqual(user, model.TbUser{}) {
			ctx.JSON(401, nil)
			fmt.Println("shibai")
			ctx.Abort()
		}
		// 有，放行
		ctx.Next()
	}
}
