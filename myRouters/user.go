package myRouters

import (
	"GoRedisLearn/controller"
	"GoRedisLearn/middleware"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
}

func (*UserRoute) InitRouter(g *gin.RouterGroup) {
	u := g.Group("/user")
	u.Use(middleware.RefreshTokenMiddleware())
	// 不需要拦截器
	{
		u.POST("/code", controller.Code)
		u.POST("/login", controller.Login)
		u.POST("/update", controller.Update)
		u.POST("/m", controller.QueryWithMutex)
		u.POST("/s", controller.QueryWithLogicExpire)
		u.POST("/blog", controller.LikeBlog)
		u.POST("/5", controller.LikeBlogTop5)
		u.POST("/sb", controller.SaveBolg)
		u.POST("/follow", controller.Follow)
		u.POST("/followOrNot", controller.IsFollow)
		u.POST("/common", controller.Common)
		u.POST("/fbr", controller.FollowByRedis)
		u.POST("/cbr", controller.CommonByRedis)
		u.POST("/qbof", controller.QueryBlogOfFollow)
		u.POST("/lsd", controller.LoadShopData)
		u.POST("/qsbt", controller.QueryShopByType)
		u.POST("/sign", controller.Sign)
		u.POST("/sc", controller.SignCount)
		u.POST("/hl", controller.HyperLog)
		u.POST("/cjb", controller.Caojinbo)

		u.POST("/xsms", controller.SecKillVoucher)   // 限时秒杀 瞎写的 自己也看不懂了
		u.POST("/xsms2", controller.SecKillVoucher2) // 限时秒杀 悲观锁
		u.POST("/xsms3", controller.SecKillVoucher3) // 限时秒杀 乐观锁
		u.POST("/xsms4", controller.SecKillVoucher4) // 限时秒杀 不加锁
		u.POST("/xsms5", controller.SecKillVoucher5) // 限时秒杀 分布式锁
		u.POST("/xsms6", controller.SecKillVoucher6) // 限时秒杀 分布式锁带看门狗
	}
	// 需要拦截器
	c := g.Group("/user")
	c.Use(middleware.RefreshTokenMiddleware(), middleware.AuthMiddleware())
	{
		c.POST("/id", controller.QueryById)
	}
}
