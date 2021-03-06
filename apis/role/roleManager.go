package role

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

func Store(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/store", func(c *gin.Context) {
		var param params.RoleParams
		err := c.BindJSON(&param)
		name := param.Name
		role := service.NewRole(conf)
		myLogger := my_log.NewLog(conf).Logger

		data, err1 := json.Marshal(param)
		if err1 != nil {
			myLogger.Info("Json marshaling failed：%s", err1)
		}
		myLogger.Info("params " + string(data))

		if err == nil {
			if name == "" {
				myLogger.Info(common.NameEmptyMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameEmptyMsg)
				return
			}

			//check duplicate name
			roleInfo := role.ShowByName(name)
			if 0 != roleInfo.Id {
				myLogger.Info(common.NameUniqueMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
				return
			}

			//insert a new record
			insertErr := role.Store(param)

			if insertErr == nil {
				myLogger.Info(common.SaveSuccessMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.SaveSuccessMsg)
			} else {
				myLogger.Error(common.SaveFailureMsg)
				apis.FormatResponseWithoutData(c, common.FailureCode, common.SaveFailureMsg)
			}
		} else {
			myLogger.Error(common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func Show(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/show/:id", func(c *gin.Context) {
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		role := service.NewRole(conf)
		myLogger := my_log.NewLog(conf).Logger

		myLogger.Info("Role info id: " + stringId)

		roleInfo := role.Show(id)
		//record does not exist
		if 0 == roleInfo.Id {
			myLogger.Info(common.RoleNotFoundMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RoleNotFoundMsg)
			return
		}

		apis.FormatResponse(c, common.SuccessCode, "", &params.Item{
			Id:        roleInfo.Id,
			Name:      roleInfo.Name,
			UpdatedAt: roleInfo.UpdatedAt,
			CreatedAt: roleInfo.CreatedAt,
		})
	})
}

func List(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/list", func(c *gin.Context) {
		role := service.NewRole(conf)
		result := role.List()
		apis.FormatResponse(c, common.SuccessCode, "", result)
	})
}

func Update(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/update/:id", func(c *gin.Context) {
		var param params.RoleParams
		idString := c.Param("id")
		id, _ := strconv.Atoi(idString)
		err := c.BindJSON(&param)
		role := service.NewRole(conf)
		myLogger := my_log.NewLog(conf).Logger

		data, err1 := json.Marshal(param)
		if err1 != nil {
			myLogger.Info("Json marshaling failed：%s", err1)
		}
		myLogger.Info("params " + string(data))

		if err == nil {
			roleInfo := role.Show(id)
			//updating role does not exit
			if 0 == roleInfo.Id {
				myLogger.Info(common.RoleNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RoleNotFoundMsg)
				return
			}

			name := param.Name

			if name == "" {
				myLogger.Info(common.NothingToUpdate)
				apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.NothingToUpdate)
				return
			} else {
				roleInfo.Name = param.Name
				roleInfo.UpdatedAt = time.Now()
			}

			//check duplicate name
			if 0 != role.UpdateShowByName(name, id).Id {
				myLogger.Info(common.NameUniqueMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
				return
			}

			updateErr := role.Update(roleInfo)
			if updateErr == nil {
				myLogger.Info(common.UpdateSuccessMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.UpdateSuccessMsg)
			} else {
				myLogger.Error(common.UpdateFailureMsg)
				apis.FormatResponseWithoutData(c, common.FailureCode, common.UpdateFailureMsg)
			}
		} else {
			myLogger.Error(common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func DeleteRole(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/delete/:id", func(c *gin.Context) {
		role := service.NewRole(conf)
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		roleInfo := role.Show(id)
		myLogger := my_log.NewLog(conf).Logger

		myLogger.Info("Delete role id： " + stringId)

		//record not found
		if 0 == roleInfo.Id {
			myLogger.Info(common.RoleNotFoundMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RoleNotFoundMsg)
			return
		}

		delError := role.Delete(roleInfo)
		if delError == nil {
			myLogger.Info(common.DeleteSuccessMsg)
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
		} else {
			myLogger.Error(delError.Error())
			apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
		}
	})
}
