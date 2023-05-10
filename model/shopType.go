package model

import "time"

type TbShopType struct {
	Id         int
	Name       string
	Icon       string
	Sort       int
	CreateTime time.Time
	UpdateTime time.Time
}
