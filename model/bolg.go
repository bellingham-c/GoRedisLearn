package model

import "time"

type TbBlog struct {
	Id         int
	ShopId     int
	UserId     int
	Title      string
	Images     string
	Content    string
	Liked      string
	Comments   string
	CreateTime time.Time
	UpdateTime time.Time
}
