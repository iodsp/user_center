package server

import (
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/apis/domain"
	"github.com/iodsp/user_center/apis/role"
	"github.com/iodsp/user_center/context"
)

func registerRoutes(app *gin.Engine, conf *context.Config) {
	//routes
	rolePrefix := app.Group("/role")
	{
		role.Store(rolePrefix, conf)
		role.Show(rolePrefix, conf)
		role.Update(rolePrefix, conf)
		role.DeleteRole(rolePrefix, conf)
		role.List(rolePrefix, conf)
	}
	domainPrefix := app.Group("/domain")
	{
		domain.Store(domainPrefix, conf)
		domain.Show(domainPrefix, conf)
		domain.List(domainPrefix, conf)
		domain.Update(domainPrefix, conf)
		domain.DeleteDomain(domainPrefix, conf)
	}
}
