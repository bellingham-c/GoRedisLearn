package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"GoRedisLearn/response"
	"GoRedisLearn/util"
	"github.com/gin-gonic/gin"
	"time"
)

func Code(c *gin.Context) {
	rds := RedisUtil.RedisUtil

	phone := c.PostForm("phone")
	//1.校验手机号
	if !util.CheckMobile(phone) {
		//2.不符合 返回错误信息
		response.Fail(c, nil)
	}
	//3.符合，生成验证码
	cacheCode := util.GenerateCode()
	//4.保存验证码到redis 5分钟过期
	//rds.Set(c, phone, cacheCode, 5*60*time.Second)
	err := rds.SET(phone, cacheCode, 5*60*time.Second)
	if err != nil {
		return
	}
	//5.发送验证码
	response.Success(c, gin.H{"data": cacheCode})
}

func Login(c *gin.Context) {
	db := DB.GetDB()
	rds := RedisUtil.RedisUtil

	phone := c.PostForm("phone")
	code := c.PostForm("code")
	//1.校验手机号
	if !util.CheckMobile(phone) {
		//2.不符合，返回错误信息
		response.Fail(c, nil)
	}
	//3.校验验证码
	cacheCode := rds.Get(c, phone).Val()
	if code == "" || cacheCode != code {
		//4.不一致，报错
		response.Fail(c, nil)
	}
	//5.一致，更具手机号查询用户
	var user model.TbUser
	tx := db.Where("phone=?", phone).Find(&user)
	//6.判断用户是否存在
	if tx.RowsAffected == 0 {
		//7.不存在，创建新用户并保存
		user.Phone = phone
		user.NickName = "user"
		db.Save(&user)
	}
	//8 保存用户信息到redis
	// 随机生成token作为登录令牌
	token, err := util.ReleaseToken(user)
	if err != nil {
		panic(err)
	}
	//8.2 将user对象转为hash存储  重写的hmset方法封装了struct to map方法
	//8.3 将用户存储到redis
	err = rds.HMSET(token, user)
	if err != nil {
		panic(err)
	}
	//8.4 设置有效期
	err = rds.EXPIRE(token, 30*time.Minute)
	if err != nil {
		panic(err)
	}
	//9 返回token
	response.Success(c, gin.H{"token": token})
}
