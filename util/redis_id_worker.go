package util

import (
	"GoRedisLearn/RedisUtil"
	"context"
	"time"
)

var (
	BEGIN_TIMESTAMP int64 = 1684472874
	ctx                   = context.Background()
)

func NextId(keyPrefix string) int64 {
	rds := RedisUtil.RedisUtil
	t := time.Now()
	// 1 生成时间戳
	NowSecond := t.Unix()
	tameStamp := NowSecond - BEGIN_TIMESTAMP
	// 2 生成序列号
	// 2.1 获取当前日期 精到天
	data := t.Format("2006:01:02")
	// 2.2 自增长
	count := rds.IncrBy(ctx, "icr:"+keyPrefix+":"+data, 1)
	// 3 拼接并返回
	return tameStamp<<32 | count.Val()
}
