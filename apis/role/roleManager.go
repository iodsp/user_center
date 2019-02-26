package role

import (
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/apis"
	"github.com/iodsp/user_center/common"
	"github.com/iodsp/user_center/models/fionaUserCenter"
	"strconv"
	"time"
)

type Params struct {
	Id   int    `json:"primary_key"`
	Name string `json:"name" binding:"required"`
}

type item struct {
	Id        int       `json:"primary_key"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"column:lastModTime"`
	CreatedAt time.Time `json:"column:createdTime"`
}

func Store(c *gin.Context) {
	var params Params
	err := c.BindJSON(&params)
	name := params.Name
	var role fionaUserCenter.Role

	if err == nil {
		if (name == "") {
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameEmptyMsg)
			return
		}

		if !common.DB.Where(&fionaUserCenter.Role{Name: name}).First(&role).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.NameUniqueMsg)
			return
		}
		createTime := time.Now()
		updateTime := time.Now()
		insertErr := common.DB.Create(&fionaUserCenter.Role{
			Name:      params.Name,
			CreatedAt: createTime,
			UpdatedAt: updateTime,
		}).Error

		if insertErr == nil {
			apis.FormatResponseWithoutData(c, common.SE_CODE, common.SaveSuccessMsg)
		} else {
			apis.FormatResponseWithoutData(c, common.FailureCode, common.SaveFailureMsg)
		}
	} else {
		apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
	}
}

func Show(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)
	var role fionaUserCenter.Role

	if common.DB.Where(&fionaUserCenter.Role{Id: id}).First(&role).RecordNotFound() {
		apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.RoleNotFoundMsg)
		return
	}

	apis.FormatResponse(c, common.SuccessCode, "", &item{
		role.Id,
		role.Name,
		role.UpdatedAt,
		role.CreatedAt,
	})
}

func RoleList(c *gin.Context) {
	var roles []fionaUserCenter.Role
	common.DB.Model(&fionaUserCenter.Role{}).Find(&roles)
	apis.FormatResponse(c, common.SuccessCode, "", &roles)
}

func Update(c *gin.Context) {
	var params Params
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	err := c.BindJSON(&params)

	if err == nil {
		var role fionaUserCenter.Role
		if common.DB.Model(fionaUserCenter.Role{}).Where("id=?", id).First(&role).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.RoleNotFoundMsg)
		}
		name := params.Name
		var tmpRole fionaUserCenter.Role

		//目前只有这一个参数
		if name == "" {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.NothingToUpdate)
			return
		} else if !common.DB.Find(&tmpRole, " name = ? AND id <> ?", name, id).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.NameUniqueMsg)
			return
		} else {
			role.Name = params.Name
			role.UpdatedAt = time.Now()
		}

		updateErr := common.DB.Save(role)
		if updateErr == nil {
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.UpdateSuccessMsg)
		} else {
			apis.FormatResponseWithoutData(c, common.FailureCode, common.UpdateSuccessMsg)
		}
	} else {
		apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
	}
}

//todo 这里可以改为更新，给删除的名字加一个时间戳或者其他的标记，这样就不会出现列表中不存在却不能添加的情况
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	var role fionaUserCenter.Role
	if common.DB.Model(fionaUserCenter.Role{}).Where("id=?", id).First(&role).RecordNotFound() {
		apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.RoleNotFoundMsg)
		return
	}

	delError := common.DB.Delete(&role).Error
	if (delError == nil) {
		apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
	} else {
		apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
	}
}
