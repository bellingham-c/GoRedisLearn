package model

import "time"

type TbVoucher struct {
	Id          int
	ShopId      int
	Title       string
	SubTitle    string
	Rules       string
	PayValue    int
	ActualValue int
	Type        int
	Status      int
	CreateTime  time.Time
	UpdateTime  time.Time
}
