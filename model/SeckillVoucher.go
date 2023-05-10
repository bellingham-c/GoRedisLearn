package model

import "time"

type TbSeckillVoucher struct {
	VoucherId  int
	Stock      int
	CreateTime time.Time
	BeginTime  time.Time
	EndTime    time.Time
	UpdateTime time.Time
}
