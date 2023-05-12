package model

import "time"

type TbUser struct {
	Id         int       `json:"Id"`
	Phone      string    `json:"Phone"`
	Password   string    `json:"Password"`
	NickName   string    `json:"NickName"`
	Icon       string    `json:"Icon"`
	CreateTime time.Time `json:"CreateTime"`
	UpdateTime time.Time `json:"UpdateTime"`
}
