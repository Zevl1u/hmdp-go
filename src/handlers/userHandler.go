package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hmdp/src/beans"
	"hmdp/src/services"
	"net/http"
)

type UserHandler struct {
}

var userService = services.UserService{}

func (uh UserHandler) VerifiedCode(c *gin.Context) {
	result := userService.VerifyCode(c)
	c.JSON(http.StatusOK, result)
}

func (uh UserHandler) Login(c *gin.Context) {
	result := userService.LoginByVerCode(c)
	c.JSON(http.StatusOK, result)
}

func (uh UserHandler) Me(c *gin.Context) {
	user := sessions.Default(c).Get("userDTO").(beans.UserDTO)
	c.JSON(200, beans.Result{Success: true, Data: user})
}
