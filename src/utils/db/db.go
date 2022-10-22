package db

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error
var RedisCli *redis.Client

func init() {
	dsn := "root:root@tcp(192.168.159.128:3306)/hmdp?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	RedisCli = redis.NewClient(&redis.Options{
		Addr:     "192.168.159.128:6379",
		Password: "root",
		DB:       1,
	})
}
