package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"hmdp/src/beans"
	"hmdp/src/utils"
	"hmdp/src/utils/db"
	"time"
)

type ShopService struct {
}

// QueryShopByIdWithLogicExpire 使用逻辑过期来解决缓存击穿
func (ss ShopService) QueryShopByIdWithLogicExpire(id uint) beans.Result {
	var shop beans.Shop
	var getShopById utils.DBQueryFunc = func(id, shopPtr interface{}) error {
		err := db.DB.First(shopPtr, "id = ?", id).Error
		return err
	}
	// 从缓存中查询
	utils.QueryWithLogicalExpire(utils.CACHE_SHOP_PREFIX, id, &shop, getShopById, 10*time.Second)
	if shop == (beans.Shop{}) {
		return beans.Result{Success: false, ErrMsg: "店铺不存在！"}
	}
	return beans.Result{Success: true, Data: shop}
}

// QueryShopByIdWithMutex 根据id查询店铺信息 使用互斥锁解决缓存击穿
func (ss ShopService) QueryShopByIdWithMutex(id uint) (result beans.Result) {
	var ctx = context.Background()
	var key = fmt.Sprintf("%s%d", utils.CACHE_SHOP_PREFIX, id)
	var shop beans.Shop
	var stringCmdPtr = db.RedisCli.Get(ctx, key)
	// GetAndStore 处理在缓存中存在key的情况的函数
	// 主要实现对空值和真实商店数据的逻辑处理
	var GetAndStore = func(jsonStr string, shop *beans.Shop) beans.Result {
		if jsonStr == "" {
			return beans.Result{Success: false, ErrMsg: "不存在的商家！"}
		} else {
			if err := json.Unmarshal([]byte(jsonStr), shop); err != nil {
				panic(err)
			}
			return beans.Result{Success: true, Data: *shop}
		}
	}
	// 缓存里查到了key对应数据
	if jsonStr, err := stringCmdPtr.Result(); err == nil {
		result = GetAndStore(jsonStr, &shop)
	} else if err == redis.Nil {
		lockKey := fmt.Sprintf("%s%d", utils.MUTEX_SHOP_PREFIX, id)
		// 获取锁失败 过一段时间再尝试访问
		if !utils.TryLock(lockKey) {
			time.Sleep(100 * time.Millisecond)
			return ss.QueryShopByIdWithMutex(id)
		} else { // 获取到锁
			defer utils.Unlock(lockKey) // 方法结束释放
			// 再次检查缓存里是否能查到key对应数据
			if jsonStr, err := stringCmdPtr.Result(); err == nil {
				result = GetAndStore(jsonStr, &shop)
			} else if err == redis.Nil { // 缓存里不存在
				// 查询数据库
				if res := db.DB.First(&shop, "id = ?", id); res.Error != nil {
					// 如果记录没找到
					if res.Error == gorm.ErrRecordNotFound {
						// 存入空字符串 用以解决缓存穿透问题
						db.RedisCli.Set(ctx, key, "", time.Minute)
						result = beans.Result{Success: false, ErrMsg: "不存在的商家!"}
					} else {
						panic(res.Error)
					}
				} else { // 数据库里存在 存入redis中
					result = beans.Result{Success: true, Data: shop}
					if err = db.RedisCli.Set(ctx, key, &shop, utils.CACHE_SHOP_INFO_TTL).Err(); err != nil {
						panic(err)
					}
				}
			}
		}
	} else if err != nil {
		panic(err)
	}
	return
}

// GetAllShopType 获取所有店铺类型列表
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

// UpdateShopInfo 更新店铺信息 删除缓存 下次请求再重建缓存
func (ss ShopService) UpdateShopInfo(shop beans.Shop) beans.Result {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", utils.CACHE_SHOP_PREFIX, shop.Id)
	// 更新数据
	affectedRows := db.DB.Model(&beans.Shop{Id: shop.Id}).Updates(shop).RowsAffected
	if affectedRows < 1 {
		return beans.Result{Success: false, ErrMsg: "更新失败，请检查id"}
	}
	// 删除缓存
	if err := db.RedisCli.Del(ctx, key).Err(); err != nil {
		panic(err)
	}
	return beans.Result{Success: true}
}
