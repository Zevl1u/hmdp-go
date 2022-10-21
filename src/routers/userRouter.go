package routers

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/handlers"
	"hmdp/src/middlewares"
)

func UserRouterInit(r *gin.Engine) {
	// 创建userHandler对象
	userHandler := handlers.UserHandler{}
	// 分组路由
	userRouter := r.Group("/user")
	{
		userRouter.GET("/code", userHandler.VerifiedCode)
		userRouter.POST("/login", userHandler.Login)
		userRouter.GET("/me", middlewares.LoginInterceptor, userHandler.Me)
	}
}
