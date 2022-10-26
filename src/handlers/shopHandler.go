package handlers

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/beans"
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

func (sh ShopHandler) UpdateShopInfo(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		c.String(http.StatusOK, "错误店铺id值!")
	}
	var shopInfo beans.Shop
	err := c.ShouldBind(&shopInfo)
	if err != nil {
		panic(err)
	}
	result := shopService.UpdateShopInfo(shopInfo)
	c.JSON(http.StatusOK, result)
}
