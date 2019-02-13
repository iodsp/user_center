package domain

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/iodsp/user_center/common"
	"github.com/iodsp/user_center/models/fionaUserCenter"
	"github.com/iodsp/user_center/apis"
	"strconv"
)

type Params struct {
	Id   int    `json:"primary_key"`
	Name string `json:"name" binding:"required"`
	Type int    `json:"type"`
}

type UpdateParams struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

type item struct {
	Id        int       `json:"primary_key"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	UpdatedAt time.Time `json:"lastModTime"`
	CreatedAt time.Time `json:"createdTime"`
}

func Store(c *gin.Context) {
	var params Params
	err := c.BindJSON(&params);
	name := params.Name
	domainType := params.Type
	var domain fionaUserCenter.Domain

	if err == nil {
		if (name == "") {
			apis.FormatResponseWithoutData(c, common.PARAM_ERROR_CODE, common.NAME_EMPTY_MSG)
			return
		}

		if domainType != 1 && domainType != 2 {
			apis.FormatResponseWithoutData(c, common.PARAM_ERROR_CODE, common.DOMAIN_TYPE_ERROR_MSG)
			return
		}

		if !common.DB.Where(&fionaUserCenter.Domain{Name: name}).First(&domain).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.NAME_UNIQUE_MSG)
			return
		}
		createTime := time.Now()
		updateTime := time.Now()
		insertErr := common.DB.Create(&fionaUserCenter.Domain{
			Name:      params.Name,
			Type:      params.Type,
			CreatedAt: createTime,
			UpdatedAt: updateTime,
		}).Error

		if insertErr == nil {
			apis.FormatResponseWithoutData(c, common.SUCCESS_CODE, common.STORE_DOMAIN_SUCCESS)
		} else {
			apis.FormatResponseWithoutData(c, common.FAILURE_CODE, common.STORE_DOMAIN_FAILURE)
		}
	} else {
		apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.PARSE_PARAM_ERROR_MSG)
	}
}

func Show(c *gin.Context) {
	var domain fionaUserCenter.Domain
	id := c.Param("id")
	if common.DB.Model(&fionaUserCenter.Domain{}).Where("id = ?", id).First(&domain).RecordNotFound() {
		apis.FormatResponseWithoutData(c, common.PARAM_ERROR_CODE, common.DOMAIN_NOT_EXITS_MSG)
		return
	}
	apis.FormatResponse(c, common.SUCCESS_CODE, "", &item{
		domain.Id,
		domain.Name,
		domain.Type,
		domain.UpdatedAt,
		domain.CreatedAt,
	})
}

func DomainList(c *gin.Context) {
	var domain = []fionaUserCenter.Domain{}
	common.DB.Find(&domain)
	apis.FormatResponse(c, common.SUCCESS_CODE, "", domain)
}

func Update(c *gin.Context) {
	var params UpdateParams
	var domain fionaUserCenter.Domain
	err := c.BindJSON(&params);
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	name := params.Name
	domainType := params.Type

	if (err == nil) {
		if common.DB.Model(&fionaUserCenter.Domain{}).Where("id = ?", id).First(&domain).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.PARAM_ERROR_CODE, common.DOMAIN_NOT_EXITS_MSG)
			return
		}

		if name != "" {
			var tmpDomain fionaUserCenter.Domain
			if !common.DB.Find(&tmpDomain, " name = ? AND id <> ?", name, id).RecordNotFound() {
				apis.FormatResponseWithoutData(c, common.PARAM_ERROR_CODE, common.NAME_UNIQUE_MSG)
				return
			}
			domain.Name = name
		} else if (domainType == 0 && name == "") {
			apis.FormatResponseWithoutData(c, common.PARAM_ERROR_CODE, common.NOTHING_TO_UPDATE)
			return
		}

		if domainType != 0 {
			domain.Type = domainType
		}

		if domainType != 1 && domainType != 2 {
			apis.FormatResponseWithoutData(c, common.PARAM_ERROR_CODE, common.DOMAIN_TYPE_ERROR_MSG)
			return
		}

		domain.UpdatedAt = time.Now()
		upError := common.DB.Save(domain).Error
		if upError != nil {
			c.JSON(200, gin.H{"er": upError})
			return
			apis.FormatResponseWithoutData(c, common.FAILURE_CODE, common.UPDATE_FAILURE_MSG)
			return
		} else {
			apis.FormatResponseWithoutData(c, common.FAILURE_CODE, common.UPDATE_SUCCESS_MSG)
			return
		}
	} else {
		apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.PARSE_PARAM_ERROR_MSG)
	}
}

func DeleteDomain(c *gin.Context) {
	id := c.Param("id")
	var domain fionaUserCenter.Domain
	if common.DB.Model(fionaUserCenter.Role{}).Where("id=?", id).First(&domain).RecordNotFound() {
		apis.FormatResponseWithoutData(c, common.PARSE_PARAM_ERROR_CODE, common.ROLE_NOT_FOUND_MSG)
		return
	}

	delError := common.DB.Delete(&domain).Error
	if (delError == nil) {
		apis.FormatResponseWithoutData(c, common.SUCCESS_CODE, common.DELETE_SUCCESS_MSG)
	} else {
		apis.FormatResponseWithoutData(c, common.FAILURE_CODE, common.DELETE_FAILURE_MSG)
	}
}
