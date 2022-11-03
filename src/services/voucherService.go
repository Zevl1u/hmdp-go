package services

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

	// 获取锁
	uuthreadid, _ := uuid.NewUUID()
	userLockKey := utils.USER_LOCK + strconv.Itoa(dto.Id)
	LockIsGot, err := db.RedisCli.SetNX(c.Request.Context(), userLockKey, uuthreadid.String(), utils.USER_LOCK_TTL).Result()
	if err != nil {
		panic(err)
	}
	// 释放锁的时候，要判断当前锁是否是当前线程写入的
	// 避免一个线程执行时间过长或者阻塞 导致锁过期 其他线程获取了锁
	// 但是在释放锁时候未做判断 导致释放了别的线程的锁
	//defer func() {
	//	tid, err := db.RedisCli.Get(c.Request.Context(), userLockKey).Result()
	//	if err != nil {
	//		panic(err)
	//	}
	//	if tid == uuthreadid.String() { // 这里锁判断和释放不是原子性操作 依旧有可能出问题
	//		db.RedisCli.Del(c.Request.Context(), userLockKey)
	//	}
	//}()
	//使用lua脚本保证redis操作的原子性
	defer func() {
		var script = redis.NewScript(`
local id = redis.call('get', KEYS[1])
if(id == ARGV[1]) then
	return redis.call('del', KEYS[1])
end
return 0
`)
		keys := []string{userLockKey}
		values := []interface{}{uuthreadid.String()}
		err := script.Run(c.Request.Context(), db.RedisCli, keys, values...).Err()
		if err != nil {
			panic(err)
		}
	}()

	// 一人一单
	if !LockIsGot { // 非阻塞式锁
		return beans.Result{Success: false, ErrMsg: "一人限购一张！"}
	}
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
