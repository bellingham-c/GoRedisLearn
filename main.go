package main

import (
	"GoRedisLearn/RedisUtil"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisUtil.GetAddr(),
		Password: RedisUtil.GetPwd(), // no password set
		DB:       RedisUtil.GetDb(),  // use default DB
	})

	ctx := context.Background()

	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "height").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)

}
