package myRouters

import (
	"GoRedisLearn/controller"
	"github.com/gin-gonic/gin"
)

type LoginRoute struct {
}

func (*LoginRoute) InitRouter(g *gin.RouterGroup) {
	u := g.Group("/login")
	{
		u.POST("/code", controller.Code)
		u.POST("/login", controller.Login)
	}
}
