package handlers

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/services"
	"net/http"
)

type VoucherHandler struct {
}

var voucherService services.VoucherService

func (vh VoucherHandler) SecKillVoucher(c *gin.Context) {

	res := voucherService.SecKillVoucher(c)
	if res.Success {
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, res.ErrMsg)
	}
}
