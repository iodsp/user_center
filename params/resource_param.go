package params

import "time"

//resource input param
type ResourceParams struct {
	Id         int    `json:"Id"`
	DomainId   int    `json:"domainId" binding:"required"`
	DomainName string `json:"domainName"`
	Name       string `json:"name" binding:"required"`
	Url        string `json:"url" binding:"required"`
	Desc       string `json:"desc"`
}

//update input param
type ResourceUpdateParams struct {
	Id         int    `json:"Id"`
	DomainId   int    `json:"domainId"`
	DomainName string `json:"domainName"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	Desc       string `json:"desc"`
}

//return item
type ResourceItem struct {
	Id         int       `json:"Id"`
	DomainId   int       `json:"DomainId"`
	DomainName string    `json:"DomainName"`
	Name       string    `json:"Dame"`
	Url        string    `json:"Url"`
	Desc       string    `json:"Desc"`
	CreatedAt  time.Time `json:"CreatedTime"`
	UpdatedAt  time.Time `json:"LastModTime"`
}
