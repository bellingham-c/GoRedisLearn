package router

import (
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/myRouters"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	router := gin.New()

	myRouter := new(myRouters.SystemGroup)

	// 配置跨域
	//router.Use(middleware.CORSMiddleware())
	groupRegistry := router.Group("/")
	{
		// 使用指定拦截器
		groupRegistry.Use(PreHandle)
		myRouter.UserRoute.InitRouter(groupRegistry)
	}
	login := router.Group("/")
	{
		myRouter.LoginRoute.InitRouter(login)
	}
	return router
}

// PreHandle 请求前拦截器
func PreHandle(c *gin.Context) {
	rds := RedisUtil.RedisUtil

	fmt.Println("拦截器")
	//1.获取请求中token
	token := c.GetHeader("Authorization")
	//不存在，拦截，返回401
	if token == "" {
		c.JSON(401, nil)
		c.Abort()
	}
	//2.基于token获取redis中的用户
	result, err := rds.GET(token)
	if err != nil {
		panic(err)
	}
	fmt.Println("res", result)
	//3.判断用户是否存在
	//4.不存在 拦截 返回401
	// 放行
	c.Next()
}
