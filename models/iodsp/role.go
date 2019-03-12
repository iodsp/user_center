package iodsp

import "time"

type Role struct {
	Id        int       `gorm:"primary_key"`
	Name      string    `gorm:"name"`
	DeletedAt interface{} `gorm:"column:deletedAt"`
	CreatedAt time.Time `gorm:"column:createdTime"`
	UpdatedAt time.Time `gorm:"column:lastModTime"`
}

func (Role) TableName() string {
	return "role"
}

