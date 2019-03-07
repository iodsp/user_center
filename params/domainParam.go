package params

import "time"

type DomainParams struct {
	Id   int    `json:"primary_key"`
	Name string `json:"name" binding:"required"`
	Type int    `json:"type"`
}

type UpdateDomainParams struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

type DomainItem struct {
	Id        int       `json:"primary_key"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	UpdatedAt time.Time `json:"lastModTime"`
	CreatedAt time.Time `json:"createdTime"`
}
