package iodsp

import "time"

type RoleResource struct {
	Id          int         `gorm:"primary_key"`
	RoleId      int         `gorm:"column:roleId"`
	ResourceId  int         `gorm:"column:resourceId"`
	ResourceUrl string      `gorm:"column:resourceUrl"`
	DeletedAt   interface{} `gorm:"column:deletedAt"`
	CreatedAt   time.Time   `gorm:"column:createdTime"`
	UpdatedAt   time.Time   `gorm:"column:lastModTime"`
}

func (RoleResource) TableName() string {
	return "role_resource"
}
