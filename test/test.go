package main

import (
	"fmt"
	"hmdp/src/beans"
	"hmdp/src/utils"
)

func main() {
	u := beans.User{
		Id:       13,
		Phone:    "13641487318",
		Password: "abc123",
		NickName: "zev",
	}
	m := utils.Struct2Map(u)
	for k, v := range m {
		fmt.Printf("%v: val: %v, valtype: %T\n", k, v, v)
	}
	var userDto beans.User
	utils.Map2Struct(m, &userDto)
	fmt.Printf("%v\n", userDto)
	fmt.Printf("%v, %T\n", userDto.Id, userDto.Id)
	fmt.Printf("%v, %T\n", userDto.CreateTime, userDto.CreateTime)
}
