package handlers

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/services"
	"net/http"
	"strconv"
)

type ShopHandler struct {
}

var shopService = services.ShopService{}

func (sh ShopHandler) Query4ShopInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusOK, "错误店铺id!")
	}
	result := shopService.QueryShopById(uint(id))
	c.JSON(http.StatusOK, result)
}

func (sh ShopHandler) ShowAllShopType(c *gin.Context) {
	result := shopService.GetAllShopType()
	c.JSON(http.StatusOK, result)
}
