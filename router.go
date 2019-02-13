package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/apis/domain"
	"github.com/iodsp/user_center/apis/role"
)

func routers() *gin.Engine {
	r := gin.Default()

	r.POST("/domain/store", domain.Store)
	r.GET("/domain/show/:id", domain.Show)
	r.GET("/domain/list", domain.DomainList)
	r.POST("/domain/update/:id", domain.Update)
	r.POST("/domain/delete/:id", domain.DeleteDomain)

	r.POST("/role/store", role.Store)
	r.GET("/role/show/:id", role.Show)
	r.GET("/role/list", role.RoleList)
	r.POST("/role/update/:id", role.Update)
	r.POST("/role/delete/:id", role.DeleteRole)

	return r
}
