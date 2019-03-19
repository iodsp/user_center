package grant

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/apis"
	"github.com/iodsp/user_center/common"
	"github.com/iodsp/user_center/context"
	"github.com/iodsp/user_center/my_log"
	"github.com/iodsp/user_center/params"
	"github.com/iodsp/user_center/service"
	"time"
)

func AddUserRole(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/userrole/add", func(c *gin.Context) {
		grantService := service.NewGrant(conf)
		userService := service.NewUser(conf)
		roleService := service.NewRole(conf)
		myLogger := my_log.NewLog(conf).Logger
		var param params.UserRoleParam
		err := c.BindJSON(&param)

		data, err1 := json.Marshal(param)
		if err1 != nil {
			myLogger.Info("Json marshaling failed：%s", err1)
		}
		myLogger.Info("params " + string(data))

		if err == nil {
			if param.UserId == 0 {
				myLogger.Info(common.UserIdEmpty)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserIdEmpty)
				return
			}
			userInfo := userService.Show(param.UserId)
			if userInfo.Id == 0 {
				myLogger.Info(common.UserNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserNotFoundMsg)
				return
			}
			if param.RoleId == 0 {
				myLogger.Info(common.RoleIdEmpty)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RoleIdEmpty)
				return
			}
			roleInfo := roleService.Show(param.RoleId)
			if roleInfo.Id == 0 {
				myLogger.Info(common.RoleNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RoleNotFoundMsg)
				return
			}

			//if has been granted and softly deleted, update
			updateUserRole := grantService.ShowByUserIdRoleIdIncDel(param.UserId, param.RoleId)
			if updateUserRole.Id != 0 {
				updateUserRole.RoleName = roleInfo.Name
				updateUserRole.DomainId = userInfo.DomainId
				updateUserRole.DomainName = userInfo.DomainName
				updateUserRole.DeletedAt = nil
				updateUserRole.UpdatedAt = time.Now()

				updateLog, parErr := json.Marshal(updateUserRole)
				if parErr != nil {
					myLogger.Info("Json marshaling failed：%s", parErr)
				}
				myLogger.Info("has been granted or deleted, update Params " + string(updateLog))

				updateErr := grantService.UpdateUserRole(updateUserRole)
				if updateErr == nil {
					myLogger.Info("has been granted, update: " + common.GrantSuccessfullyMsg)
					apis.FormatResponseWithoutData(c, common.SuccessCode, common.GrantSuccessfullyMsg)
					return
				} else {
					myLogger.Error("has been granted, fail to update:" + updateErr.Error())
					apis.FormatResponseWithoutData(c, common.SuccessCode, common.GrantFailMsg)
					return
				}
			}

			param.RoleName = roleInfo.Name
			param.DomainId = userInfo.DomainId
			param.DomainName = userInfo.DomainName
			data1, err2 := json.Marshal(param)
			if err2 != nil {
				myLogger.Info("Json marshaling failed：%s", err2)
			}
			myLogger.Info("insertParams " + string(data1))

			insertErr := grantService.AddUserRole(param)
			if insertErr == nil {
				myLogger.Info(common.GrantSuccessfullyMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.GrantSuccessfullyMsg)
				return
			} else {
				myLogger.Error(common.GrantFailMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.GrantFailMsg)
				return
			}
		} else {
			myLogger.Error(common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
			return
		}
	})
}

func DeleteUserRole(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/userrole/delete", func(c *gin.Context) {
		var param params.UserRoleParam
		grantService := service.NewGrant(conf)
		myLogger := my_log.NewLog(conf).Logger
		userService := service.NewUser(conf)
		roleService := service.NewRole(conf)

		data, err := json.Marshal(param)
		if err != nil {
			myLogger.Info("Json marshaling failed：%s", err)
		}
		myLogger.Info("params " + string(data))

		err1 := c.BindJSON(&param)
		if err1 == nil {
			if param.UserId == 0 {
				myLogger.Info(common.UserIdEmpty)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserIdEmpty)
				return
			}
			userInfo := userService.Show(param.UserId)
			if userInfo.Id == 0 {
				myLogger.Info(common.UserNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserNotFoundMsg)
				return
			}
			if param.RoleId == 0 {
				myLogger.Info(common.RoleIdEmpty)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RoleIdEmpty)
				return
			}
			roleInfo := roleService.Show(param.RoleId)
			if roleInfo.Id == 0 {
				myLogger.Info(common.RoleNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RoleNotFoundMsg)
				return
			}

			userRoleInfo := grantService.ShowByUserIdRoleId(param.UserId, param.RoleId)
			if userRoleInfo.Id == 0 {
				myLogger.Info(common.UserNotRole)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserNotRole)
				return
			}

			delErr := grantService.DeleteUserRole(userRoleInfo)
			if delErr == nil {
				myLogger.Info(common.DeleteSuccessMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
				return
			} else {
				myLogger.Error(common.DeleteFailureMsg)
				apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
				return
			}
		} else {
			myLogger.Error(common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
			return
		}
	})
}
