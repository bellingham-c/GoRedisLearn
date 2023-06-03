package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var RCTX = context.Background()

func LikeBlog(c *gin.Context) {
	rds := RedisUtil.RedisUtil
	db := DB.GetDB()
	blogId := c.PostForm("id")

	// 1 获取登录用户 TODO 为了方便这里直接指明用户
	var user model.TbUser
	user.Id = 2
	// 2 判断当前登录用户是否已经点赞
	key := "blog:liked:" + blogId
	sismember := rds.SISMEMBER(key, user.Id)
	if !sismember {
		// 3 未点赞，可以点赞
		// 3.1 数据库点赞加一
		tx := db.Model("tb_blog").Where("id=?", blogId).UpdateColumn("liked", gorm.Expr("liked+?", 1))
		if tx.RowsAffected != 0 {
			// 更新成功
			// 3.2 保存用户到redis的set集合
			rds.SADD(key, user.Id)
		}
	} else {
		// 4 已点赞 取消点赞
		// 4.1 数据库点赞数减一
		db.Model("tb_blog").Where("id=?", blogId).UpdateColumn("liked", gorm.Expr("liked-?", 1))
		// 4.2把用户从redis的set集合移除
		rds.SRem(RCTX, key, user.Id)
	}

}
