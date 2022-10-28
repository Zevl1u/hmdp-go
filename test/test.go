package main

import (
	"context"
	"fmt"
	"hmdp/src/beans"
	"hmdp/src/utils"
	"hmdp/src/utils/db"
	"time"
)

func main() {
	var shop beans.Shop
	ctx := context.Background()
	var id = 2
	var key = fmt.Sprintf("%s%d", utils.CACHE_SHOP_PREFIX, id)

	if res := db.DB.First(&shop, "id = ?", id); res.Error == nil {
		if err := db.RedisCli.Set(ctx, key, &beans.LogicExpireShopInfo{
			ExpireTime: time.Now().Add(10 * time.Second),
			Shop:       shop,
		}, -1).Err(); err != nil {
			panic(err)
		}
	}
	fmt.Println(shop)
}
