package utils

import (
	"context"
	"fmt"
	"hmdp/src/utils/db"
	"math/rand"
	"time"
)

func RandomVerCode() string {
	rand.Seed(time.Now().UnixMilli())
	str := fmt.Sprintf("%.6f", rand.Float32())[2:]
	return str
}

func RandStr(length int) string {
	arr := make([]byte, length)
	for i := 0; i < length; i++ {
		rand.Seed(time.Now().UnixMilli() + int64(i))
		n := byte(rand.Intn(26))
		arr[i] = 'a' + n
	}
	return string(arr)
}

func RedisIdGenerate(keyPrefix string) int64 {
	ctx := context.Background()
	// 生成时间戳
	timeStamp := time.Now().Unix() - TIMESTAMP_BEGIN
	// 生成序列号
	today := time.Now().Format("2006-01-02")
	key := "incr" + keyPrefix + ":" + today
	uuid, err := db.RedisCli.Incr(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return (timeStamp << COUNT_BITS) | uuid
}
