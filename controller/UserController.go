package controller

import (
	"GoRedisLearn/response"
	"GoRedisLearn/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Code(c *gin.Context) {
	phone := "15071852992"
	session := "session"
	//1.校验手机号
	if !util.CheckMobile(phone) {
		//2.不符合，返回错误信息
		response.Fail(c, nil)
	}
	//3.符合，生成验证码
	code := util.GenerateCode()
	//4.保存验证码到session
	fmt.Println(session)
	//5.发送验证码
	fmt.Println(code)

	response.Success(c, nil)
}
