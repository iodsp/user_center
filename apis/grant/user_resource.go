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
	"strings"
)

func UserHasResource(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/hasuserresource", func(c *gin.Context) {
		var param params.UserResource
		err := c.BindJSON(&param)
		param.ResourceUrl = strings.Trim(param.ResourceUrl, " ")
		grantService := service.NewGrant(conf)
		userService := service.NewUser(conf)
		resourceService := service.NewResource(conf)
		myLogger := my_log.NewLog(conf).Logger

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
			roleIdInfo := grantService.ShowRoleIdsByUserId(param.UserId)

			//roleIds 1,2,3
			var roleIds string
			for _, v := range roleIdInfo {
				if v.RoleId != 0 {
					roleIds = roleIds + "," + strconv.Itoa(v.RoleId)
				}
			}
			roleIds = strings.Trim(roleIds, ",")

			roleData, roleErr := json.Marshal(roleIds)
			if err != nil {
				myLogger.Info("Json marshaling failed：%s", roleErr)
			}
			myLogger.Info("users' role ids: " + string(roleData))

			if param.ResourceUrl == "" {
				myLogger.Info(common.ResourceUrlEmpty)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.ResourceUrlEmpty)
				return
			}

			resourceInfo := resourceService.ShowByUrl(param.ResourceUrl)
			if resourceInfo.Id == 0 {
				myLogger.Info("resource " + common.RecordNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, "resource "+common.RecordNotFoundMsg)
				return
			}

			roleResource := grantService.HasRoleResource(roleIds, param.ResourceUrl)
			if roleResource.Id != 0 {
				myLogger.Info(common.UserHasResource)
				apis.FormatResponseWithoutData(c, common.UserHasResourceCode, common.UserHasResource)
				return
			} else {
				myLogger.Info(common.UserNotResource)
				apis.FormatResponseWithoutData(c, common.UserNotResourceCode, common.UserNotResource)
				return
			}
		} else {
			myLogger.Error(common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, err.Error())
			return
		}
	})
}
