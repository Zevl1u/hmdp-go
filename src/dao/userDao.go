package dao

import (
	"hmdp/src/beans"
	"hmdp/src/utils"
	"hmdp/src/utils/db"
	"log"
)

type UserDao struct {
}

func (u UserDao) CreateUserByPhone(phone string) beans.User {
	user := beans.User{Phone: phone, NickName: utils.USER_NICK_NAME_PREFIX + utils.RandStr(10)}
	res := db.DB.Create(&user)
	if res.Error != nil {
		log.Println(res.Error.Error())
	}
	return user
}

func (u UserDao) FindUserByPhone(phone string) beans.User {
	user := beans.User{}
	res := db.DB.First(&user, "phone = ?", phone)
	if res.Error != nil {
		log.Println(res.Error.Error())
	}
	return user
}
