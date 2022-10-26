package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"hmdp/src/beans"
	"hmdp/src/dao"
	"hmdp/src/utils"
	"hmdp/src/utils/db"
	"time"
)

type UserService struct {
}

var userDao dao.UserDao = dao.UserDao{}

func (us UserService) VerifyCode(c *gin.Context) beans.Result {
	time.Sleep(5 * time.Second)
	// 从表单中获取手机号
	phone := c.Query("phone")
	// 验证手机号是否有效 无效的话直接返回错误信息
	if isValid := utils.IsPhoneValid(phone); !isValid {
		return beans.Result{ErrMsg: "手机号格式错误！"}
	}
	// 生成验证码
	code := utils.RandomVerCode()
	// 保存到redis中
	err := db.RedisCli.Set(c.Request.Context(), utils.LOGIN_CODE_PREFIX+phone, code, utils.LOGIN_CODE_TTL).Err()
	if err != nil {
		panic(err)
	}
	// 假装发送验证码
	fmt.Println("发送验证码成功！验证码为：" + code)
	return beans.Result{Success: true}
}

func (us UserService) LoginByVerCode(c *gin.Context) beans.Result {
	// redis相关操作方法需要传入Context 我也不知道用这个对不对
	ctx := c.Request.Context()
	// 校验手机号
	phone := c.PostForm("phone")
	cacheCode, err := db.RedisCli.Get(ctx, utils.LOGIN_CODE_PREFIX+phone).Result()
	// 在redis不存在改手机号对应的键 证明发送验证码手机号和登录手机号不同
	if err == redis.Nil {
		return beans.Result{ErrMsg: "手机号错误！"}
	} else if err != nil {
		panic(err)
	}
	// 校验验证码
	code := c.PostForm("code")
	if cacheCode != code {
		return beans.Result{ErrMsg: "验证码错误！"}
	}
	// 一致 根据手机号查询用户
	user := userDao.FindUserByPhone(phone)
	// 若不存在的话，创建新用户并保存
	if user == (beans.User{}) {
		user = userDao.CreateUserByPhone(phone)
	}
	// 生产DTO对象，减少内存压力，也减少敏感信息传输到前端
	userDTO := beans.UserDTO{
		Id:       user.Id,
		NickName: user.NickName,
		Icon:     user.Icon,
	}
	// 保存用户信息到redis
	uuid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	db.RedisCli.Set(ctx, utils.LOGIN_CODE_PREFIX+uuid.String(), &userDTO, utils.LOGIN_USERDTO_TTL)
	// 将登录凭证写入头
	c.Header(utils.AUTHORIZATION, uuid.String())
	return beans.Result{Success: true}
}
