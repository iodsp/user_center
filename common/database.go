package common

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"fmt"
)

var DB *gorm.DB

func GetDb() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@tcp(192.168.33.11:3306)/FionaUserCenter?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	db.LogMode(true)
	DB = db
	return db
}
