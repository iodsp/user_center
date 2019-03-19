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
	"strconv"
	"time"
)

func AddRoleResource(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/addresourcerole", func(c *gin.Context) {
		grantService := service.NewGrant(conf)
		roleService := service.NewRole(conf)
		resourceService := service.NewResource(conf)
		myLogger := my_log.NewLog(conf).Logger
		var param params.RoleResource
		err1 := c.BindJSON(&param)

		data, err := json.Marshal(param)
		if err != nil {
			myLogger.Info("Json marshaling failed：%s", err)
		}
		myLogger.Info("params " + string(data))

		if err1 == nil {
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

			if param.ResourceId == 0 {
				myLogger.Info(common.ResourceIdEmpty)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.ResourceIdEmpty)
				return
			}
			resourceInfo := resourceService.Show(param.ResourceId)
			if resourceInfo.Id == 0 {
				stringReId := strconv.Itoa(param.ResourceId)
				myLogger.Info("resourceId " + stringReId + " " + common.RecordNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, "resourceId "+stringReId+" "+common.RecordNotFoundMsg)
				return
			}

			//if has been granted and softly deleted, update
			updateResourceRole := grantService.ShowByResourceIdRoleIdIncDel(param.ResourceId, param.RoleId)

			checkLog, checkErr := json.Marshal(updateResourceRole)
			if checkErr != nil {
				myLogger.Info("Json marshaling failed：%s", checkErr)
			}

			myLogger.Info("check if has been granted or deleted, update Params " + string(checkLog))
			myLogger.Info("has been granted or deleted, update Params  updateResourceRole.Id " + strconv.Itoa(updateResourceRole.Id))
			if updateResourceRole.Id != 0 {
				updateResourceRole.ResourceUrl = resourceInfo.Url
				updateResourceRole.DeletedAt = nil
				updateResourceRole.UpdatedAt = time.Now()

				updateLog, parErr := json.Marshal(updateResourceRole)
				if parErr != nil {
					myLogger.Info("Json marshaling failed：%s", parErr)
				}
				myLogger.Info("has been granted or deleted, update Params " + string(updateLog))

				updateErr := grantService.UpdateResourceRole(updateResourceRole)
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

			param.ResourceUrl = resourceInfo.Url
			insertErr := grantService.AddRoleResource(param)
			if insertErr == nil {
				myLogger.Info(common.GrantSuccessfullyMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.GrantSuccessfullyMsg)
				return
			} else {
				myLogger.Error(common.GrantFailMsg)
				apis.FormatResponseWithoutData(c, common.FailureCode, common.GrantFailMsg)
				return
			}
		} else {
			myLogger.Error(common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
			return
		}
	})
}

func DeleteResource(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/deleteresourcerole", func(c *gin.Context) {
		grantService := service.NewGrant(conf)
		myLogger := my_log.NewLog(conf).Logger
		var param params.RoleResource
		err1 := c.BindJSON(&param)
		roleId := param.RoleId
		resourceId := param.ResourceId
		roleService := service.NewRole(conf)
		resourceService := service.NewResource(conf)

		data, err := json.Marshal(param)
		if err != nil {
			myLogger.Info("Json marshaling failed：%s", err)
		}
		myLogger.Info("params " + string(data))

		if err1 == nil {
			if roleId == 0 {
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

			if resourceId == 0 {
				myLogger.Info(common.ResourceIdEmpty)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.ResourceIdEmpty)
				return
			}
			resourceInfo := resourceService.Show(param.ResourceId)
			if resourceInfo.Id == 0 {
				stringReId := strconv.Itoa(param.ResourceId)
				myLogger.Info("resourceId " + stringReId + " " + common.RecordNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, "resourceId "+stringReId+" "+common.RecordNotFoundMsg)
				return
			}

			roleResource := grantService.ShowByRoleIdResourceId(roleId, resourceId)
			if roleResource.Id == 0 {
				myLogger.Info(common.RoleNotResource)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RoleNotResource)
				return
			}

			delErr := grantService.DeleteRoleResource(roleResource)
			if delErr == nil {
				myLogger.Info(common.DeleteSuccessMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DeleteSuccessMsg)
				return
			} else {
				myLogger.Error(common.DeleteFailureMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DeleteFailureMsg)
				return
			}
		} else {
			myLogger.Error(common.ParseParamErrorCode)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
			return
		}
	})
}
