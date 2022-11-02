package main

import (
	"hmdp/src/utils/db"
)

func main() {
	affected := db.DB.Exec("update tb_seckill_voucher set stock = stock-1 where voucher_id = ? and stock > 0",
		2).RowsAffected
	println(affected)
}
