package model

import "time"

type TbSign struct {
	Id       int
	UserId   int
	Year     time.Time
	Month    time.Month
	Date     time.Time
	IsBackup int
}
