package router

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"GoRedisLearn/myRouters"
	"github.com/gin-gonic/gin"
	"reflect"
	"time"
)

var user model.TbUser

func Routers() *gin.Engine {
	router := gin.New()

	myRouter := new(myRouters.SystemGroup)

	// 配置跨域
	//router.Use(middleware.CORSMiddleware())
	groupRegistry := router.Group("/")
	{
		myRouter.UserRoute.InitRouter(groupRegistry)
	}

	return router
}

// PreHandle 请求前拦截器
func PreHandle(c *gin.Context) {
	// token 刷新拦截器
	RefreshTokenInterceptor(c)
	// 登录拦截器
	LoginInterceptor(c)
}

// LoginInterceptor 登录拦截器
func LoginInterceptor(c *gin.Context) {
	// 如果user为空 没有用户
	if reflect.DeepEqual(user, model.TbUser{}) {
		//拦截
		c.JSON(401, nil)
		c.Abort()
	}
	// 有用户 放行
	c.Next()
}

// RefreshTokenInterceptor 拦下所有请求
func RefreshTokenInterceptor(c *gin.Context) {
	rds := RedisUtil.RedisUtil
	db := DB.GetDB()

	//1.获取请求中token
	token := c.GetHeader("Authorization")
	//不存在，放行
	if token == "" {
		c.Next()
	}
	//2.基于token获取redis中的用户
	token = token[7:]
	result, err := rds.HGETALL(token)
	if err != nil {
		panic(err)
	}
	//3.判断用户是否存在
	tx := db.Where("id=?", result["Id"]).Find(&user)
	if tx.RowsAffected == 0 {
		//4.不存在，放行
		c.Next()
	}
	//刷新有效期
	err = rds.EXPIRE(token, 30*time.Minute)
	if err != nil {
		panic(err)
	}
	// 放行
	c.Next()
}
