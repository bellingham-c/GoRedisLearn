package model

import "time"

type TbUser struct {
	Id         int
	Phone      string
	Password   string
	NickName   string
	Icon       string
	CreateTime time.Time
	UpdateTime time.Time
}
