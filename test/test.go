package main

import (
	"hmdp/src/utils/db"
)

func main() {
	num := 0
	db.DB.Raw("select count(*) from tb_voucher_order where user_id = ? and voucher_id = ?", 1, 2).Scan(&num)
	println(num)

}
