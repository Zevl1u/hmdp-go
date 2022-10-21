package services

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hmdp/src/beans"
	"hmdp/src/dao"
	"hmdp/src/utils"
)

type UserService struct {
}

var userDao dao.UserDao = dao.UserDao{}

func (us UserService) VerifyCode(c *gin.Context) beans.Result {
	// 从上下文中获取参数键
	phone := c.Query("phone")
	// 验证手机号是否有效 无效的话直接返回错误信息
	if isValid := utils.IsPhoneValid(phone); !isValid {
		return beans.Result{ErrMsg: "手机号格式错误！"}
	}
	// 生成验证码
	code := utils.RandomVerCode()
	// 获取session对象
	session := sessions.Default(c)
	session.Set("code", code)   // 保存验证码到session
	session.Set("phone", phone) // 保存手机号到session 发送验证码后登录时候改了手机号
	session.Save()              // 必须要调用保存

	// 假装发送验证码
	fmt.Println("发送验证码成功！验证码为：" + code)
	return beans.Result{Success: true}
}

func (us UserService) LoginByVerCode(c *gin.Context) beans.Result {
	// 获取session对象
	session := sessions.Default(c)
	// 校验手机号
	phone := c.PostForm("phone")
	cachePhone := session.Get("phone").(string)
	if phone != cachePhone {
		return beans.Result{ErrMsg: "手机号错误！"}
	}
	// 校验验证码
	code := c.PostForm("code")
	cacheCode := session.Get("code").(string)
	if cacheCode != code {
		return beans.Result{ErrMsg: "验证码错误！"}
	}
	// 一致 根据手机号查询用户
	user := userDao.FindUserByPhone(phone)
	// 不存在的话，创建新用户并保存
	if user == (beans.User{}) {
		user = userDao.CreateUserByPhone(phone)
	}
	// 生产DTO对象，减少内存压力，也减少不必要信息传输到前端
	userDTO := beans.UserDTO{
		Id:       user.Id,
		NickName: user.NickName,
		Icon:     user.Icon,
	}
	// 保存用户信息到session
	session.Set("userDTO", userDTO)
	session.Save()
	return beans.Result{Success: true}
}
