package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hmdp/src/middlewares"
	"hmdp/src/routers"
)

func main() {
	// 初始化引擎 默认用了Logger和Recovery中间件
	r := gin.Default()
	// 使用全局中间件 这里是使用session中间件
	r.Use(middlewares.Session("hmdp")) // hmdp是设置存储在cookie中的键

	// 初始化User路由分组
	routers.UserRouterInit(r)

	r.GET("/ping", func(c *gin.Context) {
		sessions := sessions.Default(c)
		code := sessions.Get("code")
		c.JSON(200, gin.H{
			"code": code,
		})
	})

	// 监听并在 0.0.0.0:8090 上启动服务
	r.Run(":8090")
}
