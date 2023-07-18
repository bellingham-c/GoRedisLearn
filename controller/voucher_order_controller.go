package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/model"
	"GoRedisLearn/util"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

var lock sync.Mutex

func SecKillVoucher(c *gin.Context) {
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

	// TODO 使用go自身的悲观锁进行加锁操作
	//lock.Lock()
	//vId, err := strconv.Atoi(id)
	//if err != nil {
	//	panic(err)
	//}
	//CreateVoucherOrder(vId, c)
	//lock.Unlock()

	// TODO 使用分布式锁来进行加锁，即使用redis的setnx方法 来确保每个用户每次只有一个请求能生效
	uid := uuid.NewV4()
	util.TryLock(uid.String(), 30)
	vId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	CreateVoucherOrder(vId, c)
	util.UnLock(uid.String())
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
