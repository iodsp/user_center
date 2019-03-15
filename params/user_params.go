package params

import (
	"time"
)

//user input param
type UserParams struct {
	Id               int         `json:"id"`
	Name             string      `json:"name" binding:"required"`
	Password         string      `json:"password" binding:"required"`
	Phone            string      `json:"phone" binding:"required"`
	RegisterPlatform int         `json:"registerPlatform"`
	DomainName       string      `json:"domainName"`
	Type             int         `json:"type"`
	IdentityId       string      `json:"identityId"`
	DomainId         int         `json:"domainId" binding:"required"`
	DeletedAt        interface{} `json:"column:deletedAt"`
	CreatedAt        time.Time   `json:"column:createdTime"`
	UpdatedAt        time.Time   `json:"column:lastModTime"`
}

//user update param
type UserUpdateParams struct {
	Id               int         `json:"id"`
	Name             string      `json:"name"`
	Password         string      `json:"password"`
	Phone            string      `json:"phone"`
	RegisterPlatform int         `json:"registerPlatform"`
	DomainName       string      `json:"domainName"`
	Type             int         `json:"type"`
	IdentityId       string      `json:"identityId"`
	DomainId         int         `json:"domainId"`
	DeletedAt        interface{} `json:"column:deletedAt"`
	CreatedAt        time.Time   `json:"column:createdTime"`
	UpdatedAt        time.Time   `json:"column:lastModTime"`
}

type UserItem struct {
	Id               int       `json:"Id"`
	Name             string    `json:"Name"`
	Phone            string    `json:"Phone" `
	RegisterPlatform int       `json:"RegisterPlatform"`
	DomainName       string    `json:"DomainName"`
	Type             int       `json:"Type"`
	IdentityId       string    `json:"IdentityId"`
	CreatedAt        time.Time `json:"CreatedTime"`
	UpdatedAt        time.Time `json:"LastModTime"`
}
