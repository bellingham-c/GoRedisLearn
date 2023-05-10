package model

import "time"

type TbFollow struct {
	Id           int
	UserId       int
	FollowUserId int
	CreateTime   time.Time
}
