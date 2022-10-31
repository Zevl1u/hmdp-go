package beans

import (
	"encoding/json"
	"time"
)

type VoucherOrder struct {
	Id         uint      `json:"id,omitempty" gorm:"primaryKey"`
	UserId     uint      `json:"user_id,omitempty"`
	VoucherId  uint      `json:"voucher_id,omitempty"`
	PayType    uint8     `json:"pay_type,omitempty"`
	Status     uint8     `json:"status,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"default:null"`
	PayTime    time.Time `json:"pay_time"`
	UseTime    time.Time `json:"use_time"`
	RefundTime time.Time `json:"refund_time"`
	UpdateTime time.Time `json:"update_time" gorm:"default:null"`
}

func (vo *VoucherOrder) TableName() string {
	return "tb_voucher_order"
}

// MarshalBinary 用于实现HSet函数自动转换成json字符串存入redis
func (vo *VoucherOrder) MarshalBinary() ([]byte, error) {
	return json.Marshal(vo)
}

// UnmarshalBinary 用于json字符串转换成结构体
func (vo *VoucherOrder) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, vo)
}
