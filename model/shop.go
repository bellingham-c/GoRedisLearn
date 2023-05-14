package model

import "time"

type TbShop struct {
	Id         int       `json:"Id"`
	Name       string    `json:"Name"`
	TypeId     int       `json:"TypeId"`
	Images     string    `json:"Images"`
	Area       string    `json:"Area"`
	Address    string    `json:"Address"`
	X          float64   `json:"X"`
	Y          float64   `json:"Y"`
	AvgPrice   int       `json:"AvgPrice"`
	Sold       int       `json:"Sold"`
	Comments   int       `json:"Comments"`
	Score      int       `json:"Score"`
	OpenHours  string    `json:"OpenHours"`
	CreateTime time.Time `json:"CreateTime"`
	UpdateTime time.Time `json:"UpdateTime"`
}
