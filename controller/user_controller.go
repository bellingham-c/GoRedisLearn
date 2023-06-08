package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"GoRedisLearn/response"
	"GoRedisLearn/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
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

func Sign(c *gin.Context) {
	rds := RedisUtil.RedisUtil
	// 1 获取当前用户
	id := "1"
	// 2 获取日期
	t := ":" + time.Now().Format("200601")
	// 3 拼接key
	key := "sign:" + id + t
	// 4 获取今天是本月的第几天
	m := time.Now().Day()
	// 5 写入redis setbit key offset 1
	rds.SetBit(RCTX, key, int64(m-1), 1)
}

func SignCount(c *gin.Context) {
	rds := RedisUtil.RedisUtil
	// 1 获取当前用户
	id := "1"
	// 2 获取日期
	t := ":" + time.Now().Format("200601")
	// 3 拼接key
	key := "sign:" + id + t
	// 4 获取今天是本月的第几天
	m := time.Now().Day()
	// 5 获取截至今天为止所有的签到记录 返回的是十进制数字 bitfield key get u14 0
	field := rds.BitField(RCTX, key, "get", "u"+strconv.Itoa(m), 0)
	num := field.Val()[0]
	fmt.Printf("%T\n", num)
	// 6 循环遍历
	count := 0
	for true {
		// 6.1 让这个数字与1做与运算 得到数字的最后一个bit位
		// 6.2 判断这个bit位是否为0
		if num&1 == 0 {
			// 6.3 为0 未签到结束
			break
		} else {
			// 6.4 不为0 已签到 计数器加一
			count++
		}
		// 6.5 把数字右移一位 抛弃最后一个bit位 继续下一个bit位
		num >>= 1
	}
	fmt.Println("count:", count)
}

func HyperLog(c *gin.Context) {
	rds := RedisUtil.RedisUtil
	//var values [1000]string
	//j := 0
	//for i := 0; i < 1000000; i++ {
	//	j = i % 1000
	//	values[j] = "user_" + strconv.Itoa(i)
	//	if j == 999 {
	//		value, _ := json.Marshal(values)
	//		add := rds.PFAdd(RCTX, "hl2", value)
	//		if add.Err() != nil {
	//			panic(add.Err())
	//		}
	//	}
	//}

	// TODO 这里有并发安全
	for i := 0; i < 1000000; i++ {
		rds.PFAdd(RCTX, "hl2", i)
	}

	fmt.Println("values:", rds.PFCount(RCTX, "hl2").Val())
	// res 1009972
}
