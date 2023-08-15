package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/model"
	"GoRedisLearn/newRedisLock"
	"GoRedisLearn/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// 不加锁
func SecKillVoucher4(c *gin.Context) {
	fmt.Println(4)
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
	// 5 扣减库存
	temp := sv.Stock - 1
	tx := db.Model(&sv).Where("voucher_id=?", id).Update("stock", temp)
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
	vId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	vo.VoucherId = vId

	// time
	vo.CreateTime = time.Now()
	vo.UpdateTime = time.Now()
	vo.PayTime = time.Now()
	vo.UseTime = time.Now()
	vo.RefundTime = time.Now()
	// 保存订单
	tx = db.Save(&vo)
	if tx.RowsAffected == 0 {
		c.JSON(500, gin.H{"msg": "库存不足"})
		return
	}
	// 7 返回订单id
	c.JSON(200, gin.H{"data": vo.Id, "msg": "success"})
}

// 使用悲观锁直接锁死进程
func SecKillVoucher2(c *gin.Context) {
	lock.Lock()
	db := DB.GetDB()
	// 1 查寻优惠券
	id := c.PostForm("id")
	var sv model.TbSeckillVoucher
	db.Where("voucher_id=?", id).Find(&sv)
	fmt.Println(sv)
	// 4 判断库存是否充足
	if sv.Stock <= 0 {
		// 库存不足 返回错误信息
		c.JSON(500, gin.H{"msg": "库存不足"})
		return
	}

	// TODO 使用go自身的悲观锁进行加锁操作
	vId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	createOrder(c, vId)
	lock.Unlock()
}

// 乐观锁解决超卖问题
func SecKillVoucher3(c *gin.Context) {
	db := DB.GetDB()
	// 1 查寻优惠券
	id := c.PostForm("id")
	var sv model.TbSeckillVoucher
	db.Where("voucher_id=?", id).Find(&sv)
	fmt.Println(sv)
	// 4 判断库存是否充足
	// 乐观锁解决问题
	if sv.Stock < 1 {
		// 库存不足 返回错误信息
		c.JSON(500, gin.H{"msg": "库存不足"})
		return
	}

	vId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	// 5 扣减库存
	// 使用乐观锁
	tx := db.Debug().Table("tb_seckill_voucher").Where("voucher_id=?", vId).Where("stock>0").UpdateColumn("stock", gorm.Expr("stock-?", 1))
	if tx.RowsAffected == 0 {
		// 扣减失败 返回错误信息
		fmt.Println("扣减库存失败")
		c.JSON(500, gin.H{"msg": "扣减库存失败"})
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
	vo.VoucherId = vId

	// time
	vo.CreateTime = time.Now()
	vo.UpdateTime = time.Now()
	vo.PayTime = time.Now()
	vo.UseTime = time.Now()
	vo.RefundTime = time.Now()
	// 保存订单

	tx = db.Save(&vo)
	if tx.RowsAffected == 0 {
		c.JSON(500, gin.H{"msg": "库存不足"})
		return
	}
	// 7 返回订单id
	c.JSON(200, gin.H{"data": vo.Id, "msg": "success"})
}

// 分布式锁解决超卖问题
func SecKillVoucher5(c *gin.Context) {
	db := DB.GetDB()
	// 查寻优惠券
	id := c.PostForm("id")
	var sv model.TbSeckillVoucher
	// 获取锁
	uid := uuid.NewV4()

	for {
		db.Where("voucher_id=?", id).Find(&sv)
		// 判断库存是否充足
		if sv.Stock < 1 {
			// 库存不足 返回错误信息
			c.JSON(500, gin.H{"msg": "库存不足"})
			return
		}
		if util.TryLock(uid.String(), c) {
			break
		}
	}

	vId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	createOrder(c, vId)
	// 释放锁
	util.UnLock(uid.String(), c)
}

// 带有看门狗的分布式锁
func SecKillVoucher6(c *gin.Context) {
	fmt.Println("watch dog is coming")
	db := DB.GetDB()
	// 查寻优惠券
	id := c.PostForm("id")
	var sv model.TbSeckillVoucher
	// 获取锁

	//for {
	//	db.Where("voucher_id=?", id).Find(&sv)
	//	// 判断库存是否充足
	//	if sv.Stock < 1 {
	//		// 库存不足 返回错误信息
	//		c.JSON(500, gin.H{"msg": "库存不足"})
	//		return
	//	}
	//	if util.TryLock(uid.String(), c) {
	//		break
	//	}
	//}

	r := newRedisLock.NewClient()
	rdsLock := newRedisLock.NewRedisLock("test_key", r)

	for {
		db.Where("voucher_id=?", id).Find(&sv)
		// 判断库存是否充足
		if sv.Stock < 1 {
			// 库存不足 返回错误信息
			c.JSON(500, gin.H{"msg": "库存不足"})
			return
		}
		// 上锁成功
		if rdsLock.Lock(c) == nil {
			break
		}
	}

	vId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	createOrder(c, vId)
	// 释放锁
	err = rdsLock.Unlock(c)
	if err != nil {
		errors.New("fail to release lock")
		return
	}
}

func createOrder(c *gin.Context, voucherId int) {
	db := DB.GetDB()

	// 一人一单
	//var tvo model.TbVoucherOrder
	//count := db.Where("user_id=1").Where("voucher_id=?", voucherId).Find(&tvo)
	//if count.RowsAffected > 0 {
	//	c.JSON(500, gin.H{"msg": "已经买过了"})
	//	c.Abort()
	//	return
	//}
	// 5 扣减库存
	tx := db.Table("tb_seckill_voucher").Where("voucher_id=?", voucherId).UpdateColumn("stock", gorm.Expr("stock-?", 1))
	if tx.RowsAffected == 0 {
		// 扣减失败 返回错误信息
		fmt.Println("扣减库存失败")
		c.JSON(500, gin.H{"msg": "扣减库存失败"})
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

	tx = db.Save(&vo)
	if tx.RowsAffected == 0 {
		c.JSON(500, gin.H{"msg": "订单创建失败"})
		return
	}
	// 7 返回订单id
	c.JSON(200, gin.H{"data": vo.Id, "msg": "success"})
}
