package params

import "time"

//role input param
type RoleParams struct {
	Id   int    `json:"primary_key"`
	Name string `json:"name" binding:"required"`
}

//role output param
type Item struct {
	Id        int       `json:"primary_key"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"column:lastModTime"`
	CreatedAt time.Time `json:"column:createdTime"`
}

