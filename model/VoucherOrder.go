package model

import "time"

type TbVoucherOrder struct {
	Id         int
	UserId     int
	VoucherId  int
	PayType    int
	Status     int
	CreateTime time.Time
	UpdateTime time.Time
	PayTime    time.Time
	UseTime    time.Time
	RefundTime time.Time
}
