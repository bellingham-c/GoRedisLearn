package RedisUtil

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

var rds _RedisClient

type CacheData struct {
	Data interface{}
	TTL  time.Time
}

// SetWithLogicExpire
// 1 将任意对象序列化为json并存储在string类型的key中，并且可以设置ttl过期时间(30分钟
func SetWithLogicExpire(key string, value interface{}, minute int) {
	marshal, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	rds.Set(ctx, key, marshal, time.Minute*30)
}

// Set
// 2 将任意对象序列化为json并存储在string类型的key中，并且可以设置逻辑过期时间，用于处理缓存击穿问题
func Set(key string, value interface{}, t time.Time) {
	tempDate := CacheData{
		Data: value,
		TTL:  t,
	}
	marshal, err := json.Marshal(tempDate)
	if err != nil {
		panic(err)
	}
	rds.SET(key, marshal, 0).Error()
}

// QueryWithPassThrough
// 3 根据指定的key查询缓存，并反序列化为指定类型，利用缓存空值的方式解决缓存穿透问题
func QueryWithPassThrough(keyPrefix string, id string, dbCallBack func(id string) any) any {
	key := keyPrefix + id
	// 1 从redis查询数据
	jsonRes, err := rds.GET(key)
	if err != nil {
		panic(err)
	}
	// 2 判断是否存在
	if len(jsonRes) != 0 {
		// 3 存在 直接返回
		// TODO 根据实际需求返回
		fmt.Println("data", jsonRes)
		return jsonRes
	}
	// 判断是否为空值
	if jsonRes == "" {
		return nil
	}
	// 不存在 根据id查数据库
	res := dbCallBack(id)
	if res == nil {
		// 将空值写入redis
		rds.SETWithOutJson(key, nil, time.Minute*30).Error()
	}
	return res
}

// QueryWithLogicExpire
// 4 根据指定的key查询缓存，并反序列化为指定类型，需要利用逻辑过期解决缓存击穿问题
func QueryWithLogicExpire(c *gin.Context) {
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

// TryLock 获取锁
func TryLock(key string) bool {
	nx := rds.SetNX(ctx, key, "1", 10*time.Second)
	return nx.Val()
}

// Unlock 释放锁
func Unlock(key string) {
	err := rds.DEL(key)
	if err != nil {
		panic(err)
	}
}

func SaveShop2Redis(id string, expireSeconds int) {
	db := DB.GetDB()
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
