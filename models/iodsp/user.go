package iodsp

import "time"

type User struct {
	Id               int         `gorm:"primary_key"`
	Name             string      `gorm:"name"`
	Password         string      `gorm:"password"`
	Phone            string      `gorm:"phone"`
	RegisterPlatform int         `gorm:"column:registerPlatform"`
	DomainName       string      `gorm:"column:domainName"`
	Type             int         `gorm:"type"`
	IdentityId       string      `gorm:"column:identityId"`
	DomainId         int         `gorm:"column:domainId"`
	DeletedAt        interface{} `gorm:"column:deletedAt"`
	CreatedAt        time.Time   `gorm:"column:createdTime"`
	UpdatedAt        time.Time   `gorm:"column:lastModTime"`
}

func (User) TableName() string {
	return "user"
}
