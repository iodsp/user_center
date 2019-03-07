package fionaUserCenter

import "time"

type Domain struct {
	Id        int         `gorm:"primary_key"`
	Name      string      `gorm:"name"`
	Type      int         `gorm:"type"`
	DeletedAt  interface{} `gorm:"column:deletedAt"`
	CreatedAt time.Time   `gorm:"column:createdTime"`
	UpdatedAt time.Time   `gorm:"column:lastModTime"`
}

func (Domain) TableName() string {
	return "domain"
}
