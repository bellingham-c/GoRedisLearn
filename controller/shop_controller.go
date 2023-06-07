package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"GoRedisLearn/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"time"
)

var ctx = context.Background()

func QueryById(c *gin.Context) {
	id := c.PostForm("id")
	key := "cache:shop:" + id

	db := DB.GetDB()
	rds := RedisUtil.RedisUtil
	// 1.从redis查询商铺缓存
	shop, _ := rds.GET(key)
	var e model.TbShop

	// 2.判断是否存在
	if shop != "" && shop != "null" {
		// 3.存在直接返回
		err := json.Unmarshal([]byte(shop), &e)
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{"data": e})
		return
	}
	// 命中为空 返回错误信息
	if shop == "null" {
		response.Fail(c, gin.H{"msg": "用户信息不存在"})
	}
	// 4.不存在，根据id查询数据库
	tx := db.Where("id=?", id).Find(&e)
	if tx.RowsAffected == 0 {
		// 5.不存在，返回错误
		err := rds.SET(key, nil, 3*time.Minute)
		if err != nil {
			panic(err)
		}
		response.Fail(c, nil)
		return
	}
	// 6.存在 写入redis 设置过期时间
	err := rds.SET(key, e, 30*time.Minute)
	if err != nil {
		panic(err)
	}
	// 7.返回
	response.Success(c, gin.H{"data": e})
}

func Update(c *gin.Context) {
	var shop = model.TbShop{
		Id:   2,
		Name: "caojinbo",
	}
	err := c.ShouldBind(&shop)
	if err != nil {
		panic(err)
	}

	db := DB.GetDB()
	rds := RedisUtil.RedisUtil

	// update database
	db.Debug().Updates(shop)

	// delete cache
	err = rds.DEL("cache:shop:" + strconv.Itoa(shop.Id))
	if err != nil {
		panic(err)
	}

}

// QueryWithLogicExpire 利用逻辑过期解决缓存击穿
func QueryWithLogicExpire(c *gin.Context) {
	rds := RedisUtil.RedisUtil
	// 1.从redis查询商铺
	id := c.PostForm("id")
	key := "cache:shop:" + id
	// 2.判断是否存在
	result, _ := rds.GET(key)
	fmt.Println("res", result)
	if len(result) == 0 {
		// 3.不存在 直接返回
		c.JSON(401, gin.H{"msg": "缓存中没有数据"})
		return
	}
	// 4.命中 将json反序列化未对象
	var rd model.RedisData
	err := json.Unmarshal([]byte(result), &rd)
	if err != nil {
		panic(err)
		return
	}
	shop := rd.Data
	expire := rd.ExpireSeconds
	// 5.判断是否过期
	if expire.After(time.Now()) {
		// 5.1 未过期 直接返回商铺信息
		c.JSON(200, gin.H{"data": shop, "msg": "商铺信息"})
	}
	// 5.2 已过期 需要缓存重建
	// 6 缓存重建
	// 6.1 获取互斥锁
	lockKey := "lock:shop:" + id
	isLock := TryLock(lockKey)
	// 6.2判断获取锁是否成功
	if isLock {
		// TODO 6.3 成功，开启独立线程，实现缓存重建
		var wait sync.WaitGroup
		go func() {
			defer wait.Done()
			SaveShop2Redis(id, 20)
		}()
		wait.Wait()
		// 释放锁
		Unlock(lockKey)
	}
	// 6.4 返回过期商铺信息
	c.JSON(200, gin.H{"data": shop, "msg": "商铺信息"})
}

// QueryWithMutex 互斥锁解决缓存击穿
func QueryWithMutex(c *gin.Context) {
	id := c.PostForm("id")
	key := "cache:shop:" + id

	db := DB.GetDB()
	rds := RedisUtil.RedisUtil
	// 1.从redis查询商铺缓存
	shop, _ := rds.GET(key)
	var e model.TbShop

	// 2.判断是否存在
	if shop != "" && shop != "null" {
		// 3.存在直接返回
		err := json.Unmarshal([]byte(shop), &e)
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{"data": e})
		return
	}
	// 命中为空 返回错误信息
	if shop == "null" {
		response.Fail(c, gin.H{"msg": "用户信息不存在"})
	}
	// 4 实现缓存重建
	// 4.1 获取互斥锁
	lockKey := "lock:shop:" + id
	lock := TryLock(lockKey)
	// 4.2 判断是否成功
	if !lock {
		// 4.3 失败，休眠重试
		time.Sleep(50)
		QueryWithMutex(c)
	}
	// 4.4 不存在，根据id查询数据库
	tx := db.Debug().Where("id=?", id).Find(&e)
	if tx.RowsAffected == 0 {
		// 5.不存在，返回错误
		err := rds.SET(key, nil, 3*time.Minute)
		if err != nil {
			panic(err)
		}
		response.Fail(c, nil)
		return
	}
	// 6.存在 写入redis 设置过期时间
	err := rds.SET(key, e, 30*time.Minute)
	if err != nil {
		panic(err)
	}
	// 7 释放互斥锁
	Unlock(lockKey)
	// 8.返回
	response.Success(c, gin.H{"data": e})
}

// queryWithPassThrough 缓存穿透
func queryWithPassThrough(c *gin.Context) {
	id := c.PostForm("id")
	key := "cache:shop:" + id

	db := DB.GetDB()
	rds := RedisUtil.RedisUtil
	// 1.从redis查询商铺缓存
	shop, _ := rds.GET(key)
	var e model.TbShop

	// 2.判断是否存在
	if shop != "" && shop != "null" {
		// 3.存在直接返回
		err := json.Unmarshal([]byte(shop), &e)
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{"data": e})
		return
	}
	// 命中为空 返回错误信息
	if shop == "null" {
		response.Fail(c, gin.H{"msg": "用户信息不存在"})
	}
	// 4.不存在，根据id查询数据库
	tx := db.Where("id=?", id).Find(&e)
	if tx.RowsAffected == 0 {
		// 5.不存在，返回错误
		err := rds.SET(key, nil, 3*time.Minute)
		if err != nil {
			panic(err)
		}
		response.Fail(c, nil)
		return
	}
	// 6.存在 写入redis 设置过期时间
	err := rds.SET(key, e, 30*time.Minute)
	if err != nil {
		panic(err)
	}
	// 7.返回
	response.Success(c, gin.H{"data": e})
}

// TryLock 获取锁
func TryLock(key string) bool {
	rds := RedisUtil.RedisUtil
	nx := rds.SetNX(context.Background(), key, "1", 10*time.Second)
	return nx.Val()
}

// Unlock 释放锁
func Unlock(key string) {
	rds := RedisUtil.RedisUtil
	err := rds.DEL(key)
	if err != nil {
		panic(err)
	}
}

func SaveShop2Redis(id string, expireSeconds int) {
	db := DB.GetDB()
	rds := RedisUtil.RedisUtil
	// 1.查询店铺数据
	var shop model.TbShop
	db.Where("id=?", id).Find(&shop)
	// 2.封装逻辑过期时间
	now := time.Now()
	newTime := now.Add(time.Second * 20)
	data := model.RedisData{
		Data:          shop,
		ExpireSeconds: newTime,
	}
	// 3.写入redis 序列化为json
	jsonBytes, err := json.Marshal(data)
	fmt.Println("data", jsonBytes)
	if err != nil {
		panic(err)
	}
	key := "cache:shop:" + id
	rds.SET(key, jsonBytes, 0).Error()
}

func LoadShopData(c *gin.Context) {
	db := DB.GetDB()
	rds := RedisUtil.RedisUtil

	// 1 把店铺分组 按照typeId分组 typeID一致的放在一个集合
	// 1.1 获取全部类型
	var typeList []int
	db.Debug().Raw("select distinct type_id from tb_shop").Scan(&typeList)
	// 2 分批完成写入redis
	for _, i := range typeList {
		var shopIdTemp []model.TbShop
		// 2.1 获取同类型的店铺的集合
		key := "shop:geo:" + strconv.Itoa(i)
		db.Debug().Where("type_id=?", i).Find(&shopIdTemp)
		for _, shop := range shopIdTemp {
			// 2.2 写入redis GEOADD key 经度 纬度 member
			rds.GeoAdd(RCTX, key, &redis.GeoLocation{Longitude: shop.X, Latitude: shop.Y, Name: strconv.Itoa(shop.Id)})
		}
	}
}
