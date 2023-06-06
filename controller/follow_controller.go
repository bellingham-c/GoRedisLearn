package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// Follow 关注某人
func Follow(c *gin.Context) {
	// 获取user id TODO 这里的user id应该从 token中获取 但这里为了方便 便直接赋值
	userId := 1

	followId := c.PostForm("id")
	isFollow := c.PostForm("isFollow")

	db := DB.GetDB()

	var follow model.TbFollow
	follow.FollowUserId, _ = strconv.Atoi(followId)
	follow.UserId = userId
	follow.CreateTime = time.Now()

	// 1 判断是否关注
	if isFollow == "true" {
		// 2 关注 新增数据
		tx := db.Save(&follow).Debug()
		if tx.RowsAffected == 0 {
			panic(tx.Error)
		}
	} else {
		// 3 取关 删除
		tx := db.Debug().Where("user_id=?", userId).Where("follow_user_id=?", followId).Delete(&follow)
		if tx.RowsAffected == 0 {
			panic(tx.Error)
		}
	}
}

// IsFollow 是否关注
func IsFollow(c *gin.Context) {
	userId := 1
	followId := c.PostForm("id")

	db := DB.GetDB()

	var follow model.TbFollow
	db.Debug().Where("user_id=?", userId).Where("follow_user_id=?", followId).Find(&follow)

	if follow.Id == 0 {
		println(false)
	} else {
		println(true)
	}
}

// Common 查看共同关注
func Common(c *gin.Context) {
	//followId := c.PostForm("id")
	id := 1
	db := DB.GetDB()

	//var uid []int

	var follow []model.TbFollow
	db.Debug().Where("user_id=?", id).Select("follow_user_id").Find(&follow)
	//fmt.Println(follow)
	//for _, tbFollow := range follow {
	//	fmt.Println(tbFollow.FollowUserId)
	//}

	var user []model.TbUser
	db.Debug().Raw("select * from tb_user where id in (select follow_user_id from tb_follow where user_id=?)", id).Scan(&user)
	fmt.Println(user)
}

// FollowByRedis 使用redis的set 来存储关注的用户 方便之后的共同关注功能
func FollowByRedis(c *gin.Context) {
	userId := 2
	followId := c.PostForm("id")

	// 关注
	rds := RedisUtil.RedisUtil
	rds.SADD(strconv.Itoa(userId), followId)

	// 取关
	//rds.SREM(strconv.Itoa(userId),followId)
}

func CommonByRedis(c *gin.Context) {
	rds := RedisUtil.RedisUtil
	inter := rds.SInter(RCTX, "1", "2")
	fmt.Println(inter)
}
