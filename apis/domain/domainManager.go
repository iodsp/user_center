package domain

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
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameEmptyMsg)
			return
		}

		if domainType != 1 && domainType != 2 {
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainTypeErrorMsg)
			return
		}

		if !common.DB.Where(&fionaUserCenter.Domain{Name: name}).First(&domain).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.NameUniqueMsg)
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
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.StoreDomainSuccess)
		} else {
			apis.FormatResponseWithoutData(c, common.FailureCode, common.StoreDomainFailure)
		}
	} else {
		apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
	}
}

func Show(c *gin.Context) {
	var domain fionaUserCenter.Domain
	id := c.Param("id")
	if common.DB.Model(&fionaUserCenter.Domain{}).Where("id = ?", id).First(&domain).RecordNotFound() {
		apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
		return
	}
	apis.FormatResponse(c, common.SuccessCode, "", &item{
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
	apis.FormatResponse(c, common.SuccessCode, "", domain)
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
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
			return
		}

		if name != "" {
			var tmpDomain fionaUserCenter.Domain
			if !common.DB.Find(&tmpDomain, " name = ? AND id <> ?", name, id).RecordNotFound() {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
				return
			}
			domain.Name = name
		} else if (domainType == 0 && name == "") {
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NothingToUpdate)
			return
		}

		if domainType != 0 {
			domain.Type = domainType
		}

		if domainType != 1 && domainType != 2 {
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainTypeErrorMsg)
			return
		}

		domain.UpdatedAt = time.Now()
		upError := common.DB.Save(domain).Error
		if upError != nil {
			c.JSON(200, gin.H{"er": upError})
			return
			apis.FormatResponseWithoutData(c, common.FailureCode, common.UpdateFailureMsg)
			return
		} else {
			apis.FormatResponseWithoutData(c, common.FailureCode, common.UpdateSuccessMsg)
			return
		}
	} else {
		apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
	}
}

func DeleteDomain(c *gin.Context) {
	id := c.Param("id")
	var domain fionaUserCenter.Domain
	if common.DB.Model(fionaUserCenter.Role{}).Where("id=?", id).First(&domain).RecordNotFound() {
		apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.RoleNotFoundMsg)
		return
	}

	delError := common.DB.Delete(&domain).Error
	if (delError == nil) {
		apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
	} else {
		apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
	}
}
