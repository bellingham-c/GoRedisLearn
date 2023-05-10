package model

import "time"

type TbUserInfo struct {
	UserId     int
	City       string
	Introduce  string
	Fans       int
	Followee   int
	Gender     int
	Birthday   time.Time
	Credits    int
	Level      int
	CreateTime time.Time
	UpdateTime time.Time
}
