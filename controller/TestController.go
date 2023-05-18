package controller

import (
	"GoRedisLearn/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	var t model.Test
	c.ShouldBind(&t).Error()
	fmt.Println("t", t)
}
