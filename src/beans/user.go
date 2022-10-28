package beans

import (
	"encoding/json"
	"time"
)

//func init() {
//	gob.Register(User{}) // 使用session中间件时候才需要
//	gob.Register(UserDTO{})
//}

type User struct { // 默认对应的表名是`users`
	Id         int       `json:"id,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Password   string    `json:"password,omitempty"`
	NickName   string    `json:"nick_name,omitempty"`
	Icon       string    `json:"icon,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"default:null"` // gorm配置增/改时候 结构体属性若为零值则使用的默认值
	UpdateTime time.Time `json:"update_time" gorm:"default:null"`
}

// TableName 返回对应的表名
func (u User) TableName() string {
	return "tb_user"
}

type UserDTO struct {
	Id       int    `json:"id,omitempty"`
	NickName string `json:"nick_name,omitempty"`
	Icon     string `json:"icon,omitempty"`
}

// MarshalBinary 用于实现HSet函数自动转换成json字符串存入redis
func (ud *UserDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(ud)
}

// UnmarshalBinary 用于json字符串转换成结构体
func (ud *UserDTO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, ud)
}
