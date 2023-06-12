package connect

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var Cluster *redis.ClusterClient

func init() {
	Cluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"192.168.80.129:7000",
			"192.168.80.129:7001",
			"192.168.80.129:7002",
			"192.168.80.129:8001",
			"192.168.80.129:8002",
			"192.168.80.129:8003",
		},
		DialTimeout:  100 * time.Millisecond,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
	})
	// 发送一个ping命令,测试是否通
	s := Cluster.Do(context.Background(), "ping").String()
	fmt.Println("s:", s)
}
