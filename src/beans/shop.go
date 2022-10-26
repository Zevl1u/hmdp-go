package beans

import (
	"encoding/gob"
	"time"
)

type Shop struct {
	Id         uint      `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	TypeId     uint      `json:"type_id,omitempty"`
	Images     string    `json:"images,omitempty"`
	Area       string    `json:"area,omitempty"`
	Address    string    `json:"address,omitempty"`
	X          float64   `json:"x,omitempty"`
	Y          float64   `json:"y,omitempty"`
	AvgPrice   uint      `json:"avg_price,omitempty"`
	Sold       uint      `json:"sold,omitempty"`
	Comments   uint      `json:"comments,omitempty"`
	Score      uint      `json:"score,omitempty"`
	OpenHours  string    `json:"open_hours,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"default:null"`
	UpdateTime time.Time `json:"update_time" gorm:"default:null"`
}

type ShopType struct {
	Id         uint      `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Icon       string    `json:"icon,omitempty"`
	Sort       uint      `json:"sort,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"default:null"`
	UpdateTime time.Time `json:"update_time" gorm:"default:null"`
}

func init() {
	gob.Register(Shop{})
	gob.Register(ShopType{})
}

func (s Shop) TableName() string {
	return "tb_shop"
}

func (st ShopType) TableName() string {
	return "tb_shop_type"
}
