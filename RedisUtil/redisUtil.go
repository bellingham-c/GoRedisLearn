package RedisUtil

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

var (
	RedisUtil _RedisClient
	ctx       = context.Background()
)

const (
	USERLOGIN         = "user:login:"
	USERINFO          = "user:info:"
	USERLOCATION      = "user:location"
	COMMODITYINFO     = "cmdty:info:"
	COMMODITYHISOTRY  = "cmdty:history:"
	COMMODITYCATEGORY = "cmdty:category"
	COMMODITYCOLLECT  = "cmdty:collect:"
	ORDERINFO         = "order:info:"
)

type _RedisClient struct {
	*redis.Client
}

func InitRedis() {
	config := viper.New()
	config.SetConfigName("application")
	config.AddConfigPath("./config")
	config.SetConfigType("yaml")
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
	RedisUtil = _RedisClient{
		redis.NewClient(&redis.Options{
			Addr:     config.GetString("redis.addr"),
			Password: config.GetString("redis.password"),
			DB:       config.GetInt("redis.db"),
		}),
	}

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	ctx.Value(config.GetString("redis.password"))
	_, err = RedisUtil.Ping(ctx).Result()
	if err != nil {
		panic("连接redis失败：" + err.Error())
	}
}

func (rc *_RedisClient) GET(key string) (result string, err error) {
	result, err = rc.Client.Get(ctx, key).Result()
	return
}

func (rc *_RedisClient) SET(key string, value any, expiration time.Duration) (err error) {
	jsonStr, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = rc.Client.Set(ctx, key, jsonStr, expiration).Err()
	return
}

func (rc *_RedisClient) EXPIRE(key string, expiration time.Duration) (err error) {
	err = rc.Client.Expire(ctx, key, expiration).Err()
	return
}
func (rc *_RedisClient) DEL(key string) (err error) {
	err = rc.Client.Del(ctx, key).Err()
	return
}

func (rc *_RedisClient) HGETALL(key string) (resultMap map[string]string, err error) {
	resultMap, err = rc.Client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return resultMap, nil
}

func (rc *_RedisClient) HSET(key string, value any) (err error) {
	if err = rc.Client.HSet(ctx, key, struct2map(value)).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *_RedisClient) HSETNX(key string, field string, value any) (err error) {
	if err = rc.Client.HSetNX(ctx, key, field, value).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *_RedisClient) HSETBYMAP(key string, value any) (err error) {
	if err = rc.Client.HSet(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *_RedisClient) ZADD(key string, members ...redis.Z) (err error) {
	if err = rc.Client.ZAdd(ctx, key, members...).Err(); err != nil {
		return nil
	}
	return nil
}

func (rc *_RedisClient) ZADDNX(key string, members ...redis.Z) (err error) {
	if err = rc.Client.ZAddNX(ctx, key, members...).Err(); err != nil {
		return nil
	}
	return nil
}

func (rc *_RedisClient) ZREM(key string, members ...interface{}) (err error) {
	if err = rc.Client.ZRem(ctx, key, members).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *_RedisClient) ZREVRANGEWITHSCORES(key string, start, stop int64) (result []redis.Z) {
	result, err := rc.Client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil
	}
	return result
}

func (rc *_RedisClient) HSET1(key string, field string, value any) (err error) {
	if err = rc.HSet(ctx, key, field, value).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *_RedisClient) HMSET(key string, value any) (err error) {
	valueMap := struct2map(value)
	if err = rc.HMSet(ctx, key, valueMap).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *_RedisClient) ZREVRANGE(key string, start, stop int64) (result []string) {
	result, err := rc.Client.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil
	}
	return result
}

func (rc *_RedisClient) SADD(key string, members ...interface{}) (affect int64) {
	if affect, _ = rc.SAdd(ctx, key, members).Result(); affect == 0 {
		return 0
	}
	return affect
}

func (rc *_RedisClient) SISMEMBER(key string, member interface{}) (ok bool) {
	if ok, _ = rc.SIsMember(ctx, key, member).Result(); !ok {
		return ok
	}
	return ok
}

func (rc *_RedisClient) SREM(key string, member ...interface{}) (affect int64) {
	if affect, _ = rc.SRem(ctx, key, member).Result(); affect == 0 {
		return 0
	}
	return affect
}

func struct2map(value any) map[string]interface{} {
	return structs.Map(value)
}
