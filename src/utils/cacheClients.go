package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"hmdp/src/utils/db"
	"time"
)

type DBQueryFunc func(interface{}, interface{}) error

type LogicalExpireData struct {
	ExpireTime time.Time   `json:"expire_time"`
	DataPtr    interface{} `json:"data,omitempty"`
}

func (l *LogicalExpireData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(l)
}

func (l *LogicalExpireData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, l)
}

// SetWithLogicalExpire 将对象打上逻辑过期标签后存入redis
func SetWithLogicalExpire(prefix string, key interface{}, objPtr interface{}, expireTTL time.Duration) {
	ctx := context.Background()
	uri := fmt.Sprintf("%s%v", prefix, key)
	logicalExpireData := LogicalExpireData{
		ExpireTime: time.Now().Add(expireTTL),
		DataPtr:    objPtr,
	}
	err := db.RedisCli.Set(ctx, uri, &logicalExpireData, -1).Err()
	if err != nil {
		panic(err)
	}
}

func QueryWithLogicalExpire(prefix string, key interface{}, ptr interface{}, dbQueryer DBQueryFunc, expireTTL time.Duration) {
	ctx := context.Background()
	uri := fmt.Sprintf("%s%v", prefix, key)
	stringCmdPtr := db.RedisCli.Get(ctx, uri)
	jsonStr, err := stringCmdPtr.Result()
	if err == redis.Nil { // 缓存里不存在
		err = dbQueryer(key, ptr)
		if err == gorm.ErrRecordNotFound {
			db.RedisCli.Set(ctx, uri, "", CACHE_NULL_TTL)
		} else {
			SetWithLogicalExpire(prefix, key, ptr, expireTTL)
		}
		return
	} else if err != nil {
		panic(err)
	}

	if jsonStr == "" { // 防止缓存穿透的空值
		ptr = nil
		return
	} else {
		logicalExpireData := LogicalExpireData{DataPtr: ptr}
		err = stringCmdPtr.Scan(&logicalExpireData)
		if err != nil {
			panic(err)
		}
		if time.Now().After(logicalExpireData.ExpireTime) { // 过期的话 尝试重载缓存
			fmt.Println("第一次检查过期")
			lockKey := "QueryWithLogicalExpire:lock:" + uri
			if TryLock(lockKey) { // 获取到锁
				// doublecheck 再次获取缓存检查是否过期
				err = stringCmdPtr.Scan(&logicalExpireData)
				if err != nil {
					panic(err)
				}
				if time.Now().After(logicalExpireData.ExpireTime) { // 依旧过期
					fmt.Println("获取锁后 第二次检查过期")
					go func() {
						fmt.Println("进入了新的goroutine")
						err = dbQueryer(key, ptr)
						// 意味着在缓存里查到，但是逻辑过期，在数据库里已查不到
						if err == gorm.ErrRecordNotFound {
							db.RedisCli.Set(ctx, uri, "", CACHE_NULL_TTL)
						} else {
							SetWithLogicalExpire(prefix, key, ptr, expireTTL)
						}
						Unlock(lockKey)
					}()
				}
			}
		}
	}
}

// 获取锁函数
func TryLock(key string) bool {
	ctx := context.Background()
	success, err := db.RedisCli.SetNX(ctx, key, 1, MUTEX_MAX_TTL).Result()
	if err != nil {
		panic(err)
	}
	return success
}

// 释放锁函数
func Unlock(key string) {
	ctx := context.Background()
	if err := db.RedisCli.Del(ctx, key).Err(); err != nil {
		panic(err)
	}
}
