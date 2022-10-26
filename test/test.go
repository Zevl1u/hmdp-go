package main

import (
	"fmt"
	"hmdp/src/beans"
	"hmdp/src/utils/db"
	"reflect"
)

func main() {
	var users []beans.User
	fmt.Println(users, &users)
	fmt.Println(users == nil)
	fmt.Println(reflect.TypeOf(users), reflect.ValueOf(users))
	db.DB.Find(&users)
	fmt.Println(users, &users)
}
