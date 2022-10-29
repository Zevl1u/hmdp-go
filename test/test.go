package main

import (
	"fmt"
	"hmdp/src/beans"
	"hmdp/src/utils"
	"hmdp/src/utils/db"
	"time"
)

func main() {
	shop := beans.Shop{}
	db.DB.First(&shop, "id = ?", 1)
	utils.SetWithLogicalExpire("abcabc-", "1", &shop, 10*time.Second)
	//utils.QueryWithLogicalExpire("abcabc-", "1", &shop, nil, 10*time.Second)
	fmt.Println(shop)
}
