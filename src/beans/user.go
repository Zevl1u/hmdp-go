package beans

import (
	"encoding/gob"
	"encoding/json"
	"time"
)

func init() {
	gob.Register(User{})
	gob.Register(UserDTO{})
}

type User struct { // 默认对应的表名是`users`
	Id         int       `json:"id,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Password   string    `json:"password,omitempty"`
	NickName   string    `json:"nick_name,omitempty"`
	Icon       string    `json:"icon,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"default:null"` // gorm配置增/改时候默认值
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

func (ud *UserDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(ud)
}

func (ud *UserDTO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, ud)
}
