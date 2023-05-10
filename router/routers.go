package router

import (
	"GoRedisLearn/myRouters"
	"github.com/gin-gonic/gin"
)

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
