package services

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/beans"
	"hmdp/src/utils"
	"hmdp/src/utils/db"
	"strconv"
	"time"
)

type VoucherService struct {
}

func (vs VoucherService) SecKillVoucher(c *gin.Context) beans.Result {
	value, exists := c.Get("userDTO")
	if !exists {
		return beans.Result{Success: false, ErrMsg: "请先登录！"}
	}
	dto := value.(beans.UserDTO)
	voucherId := c.Param("voucher_id")
	id, err := strconv.Atoi(voucherId)
	if err != nil {
		panic(err)
	}
	var secKillVoucher = beans.SecKillVoucher{VoucherId: id}
	affected := db.DB.First(&secKillVoucher).RowsAffected
	if affected < 1 {
		return beans.Result{Success: false, ErrMsg: "该券非秒杀券或错误的优惠券id！"}
	}
	if secKillVoucher.Stock < 1 {
		return beans.Result{Success: false, ErrMsg: "库存不足！(查询时)"}
	}
	if time.Now().Before(secKillVoucher.BeginTime) {
		return beans.Result{Success: false, ErrMsg: "秒杀尚未开始！"}
	}
	if time.Now().After(secKillVoucher.EndTime) {
		return beans.Result{Success: false, ErrMsg: "秒杀已经结束！"}
	}
	//db.DB.Model(&secKillVoucher).Where("stock = ?", secKillVoucher.Stock).Update("stock", secKillVoucher.Stock-1)
	affected = db.DB.Exec("update tb_seckill_voucher set stock = stock-1 where voucher_id = ? and stock > 0",
		voucherId).RowsAffected
	if affected < 1 {
		return beans.Result{Success: false, ErrMsg: "库存不足！(更新时)"}
	}
	orderId := utils.RedisIdGenerate(utils.VOUCHER_ORDER_PREFIX)
	var voucherOrder = beans.VoucherOrder{
		Id:        uint(orderId),
		UserId:    uint(dto.Id),
		VoucherId: uint(id),
	}
	db.DB.Save(&voucherOrder)
	return beans.Result{Success: true, Data: voucherOrder.Id}
}
