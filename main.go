package main

import (
	"GoRedisLearn/DB"
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

	r := router.Routers()

	panic(r.Run(":9090"))

}
