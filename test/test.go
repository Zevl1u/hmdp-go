package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hmdp/src/utils/db"
)

type Stu struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Profile Info   `json:"profile,omitempty"`
}

type Info struct {
	Score   int    `json:"score,omitempty"`
	Address string `json:"address,omitempty"`
}

// 忽略了许多错误处理
func main() {
	ctx := context.Background()
	stu := Stu{3, "Leo", Info{Score: 0, Address: "guangdong"}}
	bytes, _ := json.Marshal(stu)
	fmt.Println(string(bytes))
	_ = db.RedisCli.Set(ctx, "stu_10", &stu, -1).Err()
	var stu2 Stu
	_ = db.RedisCli.Get(ctx, "stu_10").Scan(&stu2)
	fmt.Println(stu2)
}

func (s *Stu) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Stu) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
