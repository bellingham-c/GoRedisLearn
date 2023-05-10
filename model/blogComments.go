package model

import "time"

type TbBlogComments struct {
	Id         int
	BlogId     int
	UserId     int
	ParentId   int
	AnswerId   int
	Content    string
	Liked      int
	Status     int
	CreateTime time.Time
	UpdateTime time.Time
}
