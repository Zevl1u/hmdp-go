package services

import (
	"context"
	"encoding/json"
	"fmt"
	"hmdp/src/beans"
	"hmdp/src/utils"
	"hmdp/src/utils/db"
)

type ShopService struct {
}

func (ss ShopService) QueryShopById(id uint) beans.Result {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", utils.CACHE_SHOP_PREFIX, id)
	// 在redis里查询
	shopMap, err := db.RedisCli.HGetAll(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	var shop beans.Shop
	// map不为空，证明在缓存里查询到了
	if len(shopMap) != 0 {
		utils.Map2Struct(shopMap, &shop)
		return beans.Result{Success: true, Data: shop}
	}
	// 若缓存里不存在，则需要去数据库查找
	res := db.DB.First(&shop, "id = ?", id)
	if res.Error != nil {
		panic(err)
	}
	// 数据库里不存在
	if shop == (beans.Shop{}) {
		return beans.Result{Success: false, ErrMsg: "不存在该商家!"}
	}
	// 数据库里存在 存入redis中
	err = db.RedisCli.HSet(ctx, key, utils.Struct2Map(shop)).Err()
	if err != nil {
		panic(err)
	}
	return beans.Result{Success: true, Data: shop}
}

func (ss ShopService) GetAllShopType() beans.Result {
	ctx := context.Background()
	key := "shop-type-list"
	// redis里查询  -1代表取到结尾
	shopTypeListStr, err := db.RedisCli.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	var shopTypeList []beans.ShopType
	var shopType beans.ShopType
	// 如果在redis里查询到了
	if len(shopTypeListStr) != 0 {
		fmt.Println("在redis中查询到商铺分类列表")
		for _, str := range shopTypeListStr {
			json.Unmarshal([]byte(str), &shopType)
			shopTypeList = append(shopTypeList, shopType)
		}
	} else {
		// 如果没查到 证明数据库数据未加载到redis中
		fmt.Println("未在redis中查询到商铺分类列表，从数据库中查询...")
		// 查询数据库多条记录，注入到shopType切片中
		db.DB.Find(&shopTypeList)
		// 将每个shopType转成json字符串存入redis列表中
		for _, v := range shopTypeList {
			bytes, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			db.RedisCli.RPush(ctx, "shop-type-list", bytes)
		}
	}
	return beans.Result{Success: true, Data: shopTypeList}
}
