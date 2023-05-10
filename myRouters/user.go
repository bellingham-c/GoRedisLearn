package myRouters

import (
	"GoRedisLearn/controller"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
}

func (*UserRoute) InitRouter(g *gin.RouterGroup) {
	u := g.Group("user")
	{
		u.POST("code", controller.Code)
	}
}
