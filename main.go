package main

import "github.com/iodsp/user_center/common"

func main() {
	r := routers()
	db := common.GetDb()
	defer db.Close()
	r.Run(":8081")
}
