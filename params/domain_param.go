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
	Id        int       `json:"Id "`
	Name      string    `json:"Name"`
	Type      int       `json:"Type"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	CreatedAt time.Time `json:"CreatedAt"`
}
