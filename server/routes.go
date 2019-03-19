package server

import (
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/apis/domain"
	"github.com/iodsp/user_center/apis/grant"
	"github.com/iodsp/user_center/apis/resource"
	"github.com/iodsp/user_center/apis/role"
	"github.com/iodsp/user_center/apis/user"
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
	userPrefix := app.Group("/user")
	{
		user.Store(userPrefix, conf)
		user.Show(userPrefix, conf)
		user.List(userPrefix, conf)
		user.UpdateUser(userPrefix, conf)
		user.DeleteUser(userPrefix, conf)
	}
	resourcePrefix := app.Group("/resource")
	{
		resource.Store(resourcePrefix, conf)
		resource.Show(resourcePrefix, conf)
		resource.List(resourcePrefix, conf)
		resource.Update(resourcePrefix, conf)
		resource.DeleteResource(resourcePrefix, conf)
	}
	grantPrefix := app.Group("/grant")
	{
		grant.AddUserRole(grantPrefix, conf)
		grant.DeleteUserRole(grantPrefix, conf)
		grant.AddRoleResource(grantPrefix, conf)
		grant.DeleteResource(grantPrefix, conf)
		grant.UserHasResource(grantPrefix, conf)
	}
}
