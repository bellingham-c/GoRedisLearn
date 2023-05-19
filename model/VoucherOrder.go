package model

import "time"

type TbVoucherOrder struct {
	Id         int64
	UserId     int       // 下单的用户id
	VoucherId  int       // 购买的代金券id
	PayType    int       // 支付方式
	Status     int       // 订单状态
	CreateTime time.Time // 下单时间
	UpdateTime time.Time // 支付时间
	PayTime    time.Time // 核销时间
	UseTime    time.Time // 退款时间
	RefundTime time.Time // 更新时间
}
