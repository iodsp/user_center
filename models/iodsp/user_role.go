package iodsp

import "time"

type UserRole struct {
	Id         int         `gorm:"primary_key"`
	RoleId     int         `gorm:"column:roleId"`
	UserId     int         `gorm:"column:userId"`
	RoleName   string      `gorm:"column:roleName"`
	DomainId   int         `gorm:"column:domainId"`
	DomainName string      `gorm:"column:domainName"`
	DeletedAt  interface{} `gorm:"column:deletedAt"`
	CreatedAt  time.Time   `gorm:"column:createdTime"`
	UpdatedAt  time.Time   `gorm:"column:lastModTime"`
}

func (UserRole) TableName() string {
	return "role_domain_user"
}
