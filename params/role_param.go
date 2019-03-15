package params

import "time"

//role input param
type RoleParams struct {
	Id   int    `json:"primary_key"`
	Name string `json:"name" binding:"required"`
}

//role output param
type Item struct {
	Id        int       `json:"Id"`
	Name      string    `json:"Name"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	CreatedAt time.Time `json:"CreatedAt"`
}

