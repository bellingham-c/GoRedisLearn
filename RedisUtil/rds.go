package RedisUtil

import "github.com/redis/go-redis/v9"

func GetClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.192.130:6379",
		Password: "caojinbo",
		DB:       0,
	})

	return client
}
