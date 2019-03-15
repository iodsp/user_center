package iodsp

import "time"

type Resource struct {
	Id         int         `gorm:"primary_key"`
	DomainId   int         `gorm:"column:domainId"`
	DomainName string      `gorm:"column:domainName"`
	Name       string      `gorm:"name"`
	Url        string      `gorm:"url"`
	Desc       string      `gorm:"desc"`
	DeletedAt  interface{} `gorm:"column:deletedAt"`
	CreatedAt  time.Time   `gorm:"column:createdTime"`
	UpdatedAt  time.Time   `gorm:"column:lastModTime"`
}

func (Resource) TableName() string {
	return "resource"
}
