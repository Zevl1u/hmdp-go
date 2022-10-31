package beans

import (
	"encoding/json"
	"time"
)

type Voucher struct {
	Id          int       `json:"id,omitempty" gorm:"primaryKey"`
	ShopId      int       `json:"shop_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	SubTitle    string    `json:"sub_title,omitempty"`
	Rule        string    `json:"rule,omitempty"`
	PayValue    int       `json:"pay_value,omitempty"`
	ActualValue int       `json:"actual_value,omitempty"`
	Type        int       `json:"type,omitempty"`
	Statue      int       `json:"statue,omitempty"`
	CreateTime  time.Time `json:"create_time" gorm:"default:null"`
	UpdateTime  time.Time `json:"update_time" gorm:"default:null"`
}

func (v *Voucher) TableName() string {
	return "tb_voucher"
}

// MarshalBinary 用于实现HSet函数自动转换成json字符串存入redis
func (v *Voucher) MarshalBinary() ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalBinary 用于json字符串转换成结构体
func (v *Voucher) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, v)
}

type SecKillVoucher struct {
	VoucherId  int       `json:"voucher_id,omitempty" gorm:"primaryKey"`
	Stock      int       `json:"stock,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"default:null"`
	BeginTime  time.Time `json:"begin_time"`
	EndTime    time.Time `json:"end_time"`
	UpdateTime time.Time `json:"update_time" gorm:"default:null"`
}

func (skv *SecKillVoucher) TableName() string {
	return "tb_seckill_voucher"
}

// MarshalBinary 用于实现HSet函数自动转换成json字符串存入redis
func (skv *SecKillVoucher) MarshalBinary() ([]byte, error) {
	return json.Marshal(skv)
}

// UnmarshalBinary 用于json字符串转换成结构体
func (skv *SecKillVoucher) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, skv)
}
