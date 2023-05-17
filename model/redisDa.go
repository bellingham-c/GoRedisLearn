package model

import "time"

type RedisData struct {
	Data          TbShop
	ExpireSeconds time.Time
}
