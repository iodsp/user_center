package role

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/iodsp/user_center/models/fionaUserCenter"
	"github.com/iodsp/user_center/common"
	"github.com/iodsp/user_center/apis"
	"strconv"
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
			apis.FormatResponseWithoutData(c, common.PARAM_ERROR_CODE, common.NAME_EMPTY_MSG)
			return
		}

		if !common.DB.Where(&fionaUserCenter.Role{Name: name}).First(&role).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.NAME_UNIQUE_MSG)
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
			apis.FormatResponseWithoutData(c, common.SUCCESS_CODE, common.SVAE_SUCCESS_MSG)
		} else {
			apis.FormatResponseWithoutData(c, common.FAILURE_CODE, common.SAVE_FAILURE_MSG)
		}
	} else {
		apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.PARSE_PARAM_ERROR_MSG)
	}
}

func Show(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)
	var role fionaUserCenter.Role

	if common.DB.Where(&fionaUserCenter.Role{Id: id}).First(&role).RecordNotFound() {
		apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.ROLE_NOT_FOUND_MSG)
		return
	}

	apis.FormatResponse(c, common.SUCCESS_CODE, "", &item{
		role.Id,
		role.Name,
		role.UpdatedAt,
		role.CreatedAt,
	})
}

func RoleList(c *gin.Context) {
	var roles []fionaUserCenter.Role
	common.DB.Model(&fionaUserCenter.Role{}).Find(&roles)
	apis.FormatResponse(c, common.SUCCESS_CODE, "", &roles)
}

func Update(c *gin.Context) {
	var params Params
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	err := c.BindJSON(&params)

	if err == nil {
		var role fionaUserCenter.Role
		if common.DB.Model(fionaUserCenter.Role{}).Where("id=?", id).First(&role).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.ROLE_NOT_FOUND_MSG)
		}
		name := params.Name
		var tmpRole fionaUserCenter.Role

		//目前只有这一个参数
		if name == "" {
			apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.NOTHING_TO_UPDATE)
			return
		} else if !common.DB.Find(&tmpRole, " name = ? AND id <> ?", name, id).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.NAME_UNIQUE_MSG)
			return
		} else {
			role.Name = params.Name
			role.UpdatedAt = time.Now()
		}

		updateErr := common.DB.Save(role)
		if updateErr == nil {
			apis.FormatResponseWithoutData(c, common.SUCCESS_CODE, common.UPDATE_SUCCESS_MSG)
		} else {
			apis.FormatResponseWithoutData(c, common.FAILURE_CODE, common.UPDATE_SUCCESS_MSG)
		}
	} else {
		apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.PARSE_PARAM_ERROR_MSG)
	}
}

//todo 这里可以改为更新，给删除的名字加一个时间戳或者其他的标记，这样就不会出现列表中不存在却不能添加的情况
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	var role fionaUserCenter.Role
	if common.DB.Model(fionaUserCenter.Role{}).Where("id=?", id).First(&role).RecordNotFound() {
		apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.ROLE_NOT_FOUND_MSG)
		return
	}

	delError := common.DB.Delete(&role).Error
	if (delError == nil) {
		apis.FormatResponseWithoutData(c, common.SUCCESS_CODE, common.DELETE_SUCCESS_MSG)
	} else {
		apis.FormatResponseWithoutData(c, common.FAILURE_CODE, common.DELETE_FAILURE_MSG)
	}
}
