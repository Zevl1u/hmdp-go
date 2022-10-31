package beans

import (
	"encoding/json"
	"time"
)

//func init() {
//	gob.Register(Shop{})
//	gob.Register(ShopType{})
//}

type Shop struct {
	Id         uint      `json:"id,omitempty" form:"id" gorm:"primaryKey"`
	Name       string    `json:"name,omitempty" form:"name"`
	TypeId     uint      `json:"type_id,omitempty" form:"type_id"`
	Images     string    `json:"images,omitempty" form:"images"`
	Area       string    `json:"area,omitempty" form:"area"`
	Address    string    `json:"address,omitempty" form:"address"`
	X          float64   `json:"x,omitempty" form:"x"`
	Y          float64   `json:"y,omitempty" form:"y"`
	AvgPrice   uint      `json:"avg_price,omitempty" form:"avg_price"`
	Sold       uint      `json:"sold,omitempty" form:"sold"`
	Comments   uint      `json:"comments,omitempty" form:"comments"`
	Score      uint      `json:"score,omitempty" form:"score"`
	OpenHours  string    `json:"open_hours,omitempty" form:"open_hours"`
	CreateTime time.Time `json:"create_time" gorm:"default:null" form:"create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"default:null" form:"update_time"`
}

func (s *Shop) TableName() string {
	return "tb_shop"
}

func (s *Shop) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Shop) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

type ShopType struct {
	Id         uint      `json:"id,omitempty" gorm:"primaryKey"`
	Name       string    `json:"name,omitempty"`
	Icon       string    `json:"icon,omitempty"`
	Sort       uint      `json:"sort,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"default:null"`
	UpdateTime time.Time `json:"update_time" gorm:"default:null"`
}

func (st *ShopType) TableName() string {
	return "tb_shop_type"
}

func (st *ShopType) MarshalBinary() ([]byte, error) {
	return json.Marshal(st)
}

func (st *ShopType) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, st)
}

type LogicExpireShopInfo struct {
	ExpireTime time.Time `json:"expire_time"`
	Shop       Shop      `json:"shop"`
}

func (lesi *LogicExpireShopInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(lesi)
}

func (lesi *LogicExpireShopInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, lesi)
}
