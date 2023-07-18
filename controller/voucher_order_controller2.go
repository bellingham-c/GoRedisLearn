package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/model"
	"GoRedisLearn/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func SecKillVoucher2(c *gin.Context) {
	db := DB.GetDB()
	// 1 查寻优惠券
	id := c.PostForm("id")
	fmt.Println("id", id)
	var sv model.TbSeckillVoucher
	db.Where("voucher_id=?", id).Find(&sv)

	// 4 判断库存是否充足
	if sv.Stock <= 0 {
		// 库存不足 返回错误信息
		c.JSON(500, gin.H{"msg": "库存不足"})
		return
	}

	vId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	createOrder(c, vId)

}

func createOrder(c *gin.Context, voucherId int) {
	db := DB.GetDB()

	// 5 扣减库存
	tx := db.Table("tb_seckill_voucher").Where("voucher_id=?", voucherId).UpdateColumn("stock", gorm.Expr("stock-?", 1))
	if tx.RowsAffected == 0 {
		// 扣减失败 返回错误信息
		c.JSON(500, gin.H{"msg": "不能修改"})
		return
	}

	// 6 创建订单
	var vo model.TbVoucherOrder
	// 6.1 订单id
	vo.Id = util.NextId("order")
	// 6.2 用户id
	// TODO 理论上这里应该从redis中获取token，然后获取到用户id，但这里为了方便就直接自定义id了
	vo.UserId = 1
	// 6.3 代金券id
	vo.VoucherId = voucherId

	// time
	vo.CreateTime = time.Now()
	vo.UpdateTime = time.Now()
	vo.PayTime = time.Now()
	vo.UseTime = time.Now()
	vo.RefundTime = time.Now()
	// 保存订单
	fmt.Println("vo", vo)

	tx = db.Save(&vo)
	if tx.RowsAffected == 0 {
		c.JSON(500, gin.H{"msg": "库存不足"})
		return
	}
	// 7 返回订单id
	c.JSON(200, gin.H{"data": vo.Id, "msg": "success"})
}
