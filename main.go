package main

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/router"
)

func main() {
	// 初始化数据库
	db := DB.InitDB() //初始化数据库
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to close database" + err.Error())
		}
		err = sqlDB.Close()
		if err != nil {
			return
		}
	}()
	// 初始化redis
	RedisUtil.InitRedis()
	// 初始化router
	r := router.Routers()

	// 设置监听端口
	panic(r.Run(":9090"))
}
