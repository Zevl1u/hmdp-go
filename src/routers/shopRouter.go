package routers

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/handlers"
	"net/http"
)

func ShopRouterInit(r *gin.Engine) {
	// 创建userHandler对象
	shopHandler := handlers.ShopHandler{}
	// 分组路由
	shopRouter := r.Group("/shop")
	{
		shopRouter.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "商户信息首页")
		})
		// 动态路由 可通过c.Param("id")获取参数
		shopRouter.GET("/:id", shopHandler.Query4ShopInfo)
		shopRouter.GET("/type-list", shopHandler.ShowAllShopType)
		shopRouter.POST("/update-info", shopHandler.UpdateShopInfo)
	}
}
