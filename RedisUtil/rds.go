package RedisUtil

import "github.com/redis/go-redis/v9"

func GetClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "43.143.241.157:6379",
		Password: "123456",
		DB:       0,
	})

	return client
}
