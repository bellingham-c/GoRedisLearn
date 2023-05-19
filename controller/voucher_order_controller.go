package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/model"
	"GoRedisLearn/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

var lock sync.Mutex

func SeckillVoucher(c *gin.Context) {
	db := DB.GetDB()
	// 1 查寻优惠券
	id := c.PostForm("id")
	var sv model.TbSeckillVoucher
	db.Where("voucher_id=?", id).Find(&sv)
	// 2 判断秒杀是否开始
	t := time.Now()
	if t.Before(sv.BeginTime) {
		// 尚未开始 返回错误信息
		c.JSON(500, gin.H{"msg": "秒杀尚未开始"})
		return
	}
	// 3 判断秒杀是否结束
	if t.After(sv.EndTime) {
		// 结束
		c.JSON(500, gin.H{"msg": "秒杀已经结束"})
		return
	}

	// 4 判断库存是否充足
	if sv.Stock <= 0 {
		// 库存不足 返回错误信息
		c.JSON(500, gin.H{"msg": "库存不足"})
		return
	}

	lock.Lock()
	vId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	CreateVoucherOrder(vId, c)
	lock.Unlock()
}

func CreateVoucherOrder(voucherId int, c *gin.Context) {
	db := DB.GetDB()
	// 一人一单
	var tvo model.TbVoucherOrder
	count := db.Where("user_id=1").Where("voucher_id=?", voucherId).Find(&tvo)
	if count.RowsAffected > 0 {
		c.JSON(500, gin.H{"msg": "已经买过了"})
		c.Abort()
		return
	}
	// 5 扣减库存
	tx := db.Table("tb_seckill_voucher").Where("voucher_id=?", voucherId).Where("stock>0").UpdateColumn("stock", gorm.Expr("stock-?", 1))
	if tx.RowsAffected == 0 {
		// 扣减失败 返回错误信息
		c.JSON(500, gin.H{"msg": "库存不足"})
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
	// 保存订单
	tx = db.Save(&vo)
	if tx.RowsAffected == 0 {
		c.JSON(500, gin.H{"msg": "库存不足"})
		return
	}
	// 7 返回订单id
	c.JSON(200, gin.H{"data": vo.Id, "msg": "success"})
}
