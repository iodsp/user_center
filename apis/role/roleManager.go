package role

import (
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/apis"
	"github.com/iodsp/user_center/common"
	"github.com/iodsp/user_center/params"
	"github.com/iodsp/user_center/context"
	"github.com/iodsp/user_center/user_center"
	"strconv"
	"time"
)

func Store(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/store", func(c *gin.Context) {
		var param params.RoleParams
		err := c.BindJSON(&param)
		name := param.Name
		role := user_center.NewRole(conf.Db())

		if err == nil {
			if name == "" {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameEmptyMsg)
				return
			}

			//todo 根据name 查找记录

			//insert a new record
			insertErr := role.Store(param)

			if insertErr == nil {
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.SaveSuccessMsg)
			} else {
				apis.FormatResponseWithoutData(c, common.FailureCode, common.SaveFailureMsg)
			}
		} else {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func Show(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/role/show/:id", func(c *gin.Context) {
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		role := user_center.NewRole(conf.Db())

		roleInfo := role.Show(id)

		//todo 没找到记录的返回报错信息
		/*if common.DB.Where(&fionaUserCenter.Role{Id: id}).First(&role).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.RoleNotFoundMsg)
			return
		}*/

		apis.FormatResponse(c, common.SuccessCode, "", &params.Item{
			Id:        roleInfo.Id,
			Name:      roleInfo.Name,
			UpdatedAt: roleInfo.UpdatedAt,
			CreatedAt: roleInfo.CreatedAt,
		})
	})
}


func List(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/role/list", func(c *gin.Context) {
		role := user_center.NewRole(conf.Db())
		result := role.List()
		apis.FormatResponse(c, common.SuccessCode, "", result)
	})
}

func Update(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/role/update/:id", func(c *gin.Context) {
		var param params.RoleParams
		idString := c.Param("id")
		id, _ := strconv.Atoi(idString)
		err := c.BindJSON(&param)
		role := user_center.NewRole(conf.Db())

		if err == nil {
			roleInfo := role.Show(id)
			//todo 判断更新的id记录是否存在
			/*if common.DB.Model(fionaUserCenter.Role{}).Where("id=?", id).First(&role).RecordNotFound() {
				apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.RoleNotFoundMsg)
			}*/
			name := param.Name

			//目前只有这一个参数
			if name == "" {
				apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.NothingToUpdate)
				return
			} else {
				roleInfo.Name = param.Name
				roleInfo.UpdatedAt = time.Now()
			}
			/*else if !common.DB.Find(&tmpRole, " name = ? AND id <> ?", name, id).RecordNotFound() {
			    todo 判断更新的名字是否与其他的已存在的记录名字重复
				apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.NameUniqueMsg)
				return
			}*/

			updateErr := role.Update(roleInfo)
			if updateErr == nil {
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.UpdateSuccessMsg)
			} else {
				apis.FormatResponseWithoutData(c, common.FailureCode, common.UpdateSuccessMsg)
			}
		} else {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func DeleteRole(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/role/delete/:id", func(c *gin.Context) {
		role := user_center.NewRole(conf.Db())
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		roleInfo := role.Show(id)

		//todo 判断要删除的id记录是否存在
		/*if common.DB.Model(fionaUserCenter.Role{}).Where("id=?", id).First(&role).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.RoleNotFoundMsg)
			return
		}*/

		delError := role.Delete(roleInfo)
		if delError == nil {
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
		} else {
			apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
		}
	})
}
