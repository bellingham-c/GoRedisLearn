package model

import "time"

type TbShop struct {
	Id         int
	Name       string
	TypeId     int
	Images     string
	Area       string
	Address    string
	X          float64
	Y          float64
	AvgPrice   int
	Sold       int
	Comments   int
	Score      int
	OpenHours  string
	CreateTime time.Time
	UpdateTime time.Time
}
