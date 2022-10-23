package main

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/middlewares"
	"hmdp/src/routers"
)

func main() {
	// 初始化引擎 默认用了Logger和Recovery中间件
	r := gin.Default()

	// 使用全局中间件 这里是使用的中间件是为了刷新token在redis中的ttl 保证有新请求之后重置过期定时器
	r.Use(middlewares.RefreshTokenInterceptor)

	// 初始化User路由分组
	routers.UserRouterInit(r)
	routers.ShopRouterInit(r)

	// 监听并在 0.0.0.0:8090 上启动服务
	r.Run(":8090")
}
