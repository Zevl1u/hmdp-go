package services

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
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
	var shop beans.Shop
	if err := db.RedisCli.Get(ctx, key).Scan(&shop); err == redis.Nil { // 若缓存里不存在，则需要去数据库查找
		if res := db.DB.First(&shop, "id = ?", id); res.Error != nil {
			panic(err)
		} else if shop == (beans.Shop{}) { // 数据库里不存在
			return beans.Result{Success: false, ErrMsg: "不存在该商家!"}
		} else { // 数据库里存在 存入redis中
			if err := db.RedisCli.Set(ctx, key, &shop, utils.CACHE_SHOP_INFO_TTL).Err(); err != nil {
				panic(err)
			}
		}
	} else if err != nil { // shop不为零值，证明在缓存里查询到了
		panic(err)
	} else { // shop不为零值，证明在缓存里查询到了
		return beans.Result{Success: true, Data: shop}
	}
	return beans.Result{Success: true, Data: shop}
}

func (ss ShopService) GetAllShopType() beans.Result {
	ctx := context.Background()
	key := "shop-type-list"
	var shopTypeList []beans.ShopType
	// redis里查询  -1代表取到结尾
	if err := db.RedisCli.LRange(ctx, key, 0, -1).ScanSlice(&shopTypeList); err != nil {
		panic(err)
	} else {
		// 如果条件成立 即在redis中没查到
		if len(shopTypeList) == 0 {
			fmt.Println("未在redis中查询到商铺分类列表，从数据库中查询...")
			// 查询数据库多条记录，注入到shopType切片中
			db.DB.Find(&shopTypeList)
			// 将每个shopType转成json字符串存入redis列表中
			for _, v := range shopTypeList {
				err = db.RedisCli.RPush(ctx, "shop-type-list", &v).Err()
				if err != nil {
					panic(err)
				}
			}
			// 尝试直接切片添加列表 失败
			//err = db.RedisCli.RPush(ctx, "shop-type-list", &shopTypeList).Err()
			//if err != nil {
			//	panic(err)
			//}
		}
	}
	return beans.Result{Success: true, Data: shopTypeList}
}

func (ss ShopService) UpdateShopInfo(shop beans.Shop) beans.Result {
	ctx := context.Background()
	if shop.Id == 0 {
		return beans.Result{Success: false, ErrMsg: "店铺id错误！"}
	}
	key := fmt.Sprintf("%s%d", utils.CACHE_SHOP_PREFIX, shop.Id)
	// 更新数据
	db.DB.Model(&shop).Updates(shop)
	// 删除缓存
	if err := db.RedisCli.Del(ctx, key).Err(); err != nil {
		panic(err)
	}
	return beans.Result{Success: true}
}
