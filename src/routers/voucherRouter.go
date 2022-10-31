package routers

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/handlers"
)

func VoucherRouterInit(r *gin.Engine) {
	voucherHandler := handlers.VoucherHandler{}
	voucherRouter := r.Group("/voucher-order")
	{
		voucherRouter.POST("/sec-kill/:voucher_id", voucherHandler.SecKillVoucher)
	}
}
