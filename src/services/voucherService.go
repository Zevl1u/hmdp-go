package services

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/beans"
	"hmdp/src/utils"
	"hmdp/src/utils/db"
	"strconv"
	"sync"
	"time"
)

type VoucherService struct {
}

var lock = new(sync.Mutex)

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

	lock.Lock()
	defer lock.Unlock()
	// 一人一单
	var num int
	db.DB.Raw("select count(*) from tb_voucher_order where user_id = ? and voucher_id = ?", dto.Id, id).Scan(&num)
	if num > 0 {
		return beans.Result{Success: false, ErrMsg: "一人限购一张！"}
	}
	// 减少库存 乐观锁
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
