package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/apis"
	"github.com/iodsp/user_center/common"
	"github.com/iodsp/user_center/params"
	"github.com/iodsp/user_center/user_center"
	"github.com/iodsp/user_center/context"
	"strconv"
	"time"
)

//r.POST("/domain/store", domain.Store)
func Store(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/store", func(c *gin.Context) {
		var param params.DomainParams
		err := c.BindJSON(&param)
		name := param.Name
		domainType := param.Type
		domainService := user_center.NewDomain(conf.Db())

		if err == nil {
			if name == "" {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameEmptyMsg)
				return
			}

			if domainType != 1 && domainType != 2 {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainTypeErrorMsg)
				return
			}

			//todo 判断插入的name是否重复
			/*if !common.DB.Where(&fionaUserCenter.Domain{Name: name}).First(&domain).RecordNotFound() {
				apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.NameUniqueMsg)
				return
			}*/

			insertErr := domainService.StoreDomain(param)
			if insertErr == nil {
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.StoreDomainSuccess)
			} else {
				apis.FormatResponseWithoutData(c, common.FailureCode, common.StoreDomainFailure)
			}
		} else {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func Show(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/domain/show/:id", func(c *gin.Context) {
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		domainService := user_center.NewDomain(conf.Db())

		//todo 详情的id记录找不到报错
		/*if common.DB.Model(&fionaUserCenter.Domain{}).Where("id = ?", id).First(&domain).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
			return
		}*/

		domain := domainService.ShowDomain(id)
		apis.FormatResponse(c, common.SuccessCode, "", &params.DomainItem{
			domain.Id,
			domain.Name,
			domain.Type,
			domain.UpdatedAt,
			domain.CreatedAt,
		})
	})
}

func List(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/domain/list", func(c *gin.Context) {
		domainService := user_center.NewDomain(conf.Db())
		domain := domainService.DomainList()
		apis.FormatResponse(c, common.SuccessCode, "", domain)
	})
}

func Update(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/domain/update/:id", func(c *gin.Context) {
		var param params.UpdateDomainParams
		err := c.BindJSON(&param);
		idString := c.Param("id")
		id, _ := strconv.Atoi(idString)
		name := param.Name
		domainType := param.Type
		domainService := user_center.NewDomain(conf.Db())

		if err == nil {
			domain := domainService.ShowDomain(id)
			//todo 更新的id记录不存在
			/*if common.DB.Model(&fionaUserCenter.Domain{}).Where("id = ?", id).First(&domain).RecordNotFound() {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
				return
			}*/

			if name != "" {
				//var tmpDomain fionaUserCenter.Domain
				//todo name不可重复
				/*if !common.DB.Find(&tmpDomain, " name = ? AND id <> ?", name, id).RecordNotFound() {
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
					return
				}*/
				domain.Name = name
			} else if domainType == 0 && name == "" {
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
			upError := domainService.UpdateDomain(domain)
			if upError != nil {
				c.JSON(200, gin.H{"er": upError})
				return
			} else {
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.UpdateSuccessMsg)
				return
			}
		} else {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func DeleteDomain(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/domain/delete/:id", func(c *gin.Context) {
		idString := c.Param("id")
		id, _ := strconv.Atoi(idString)
		domainService := user_center.NewDomain(conf.Db())

		//todo 删除的id记录不存在
		domain := domainService.ShowDomain(id)
		/*if common.DB.Model(fionaUserCenter.Role{}).Where("id=?", id).First(&domain).RecordNotFound() {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.RoleNotFoundMsg)
			return
		}*/

		delError := domainService.DeleteDomain(domain)
		if delError == nil {
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
		} else {
			apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
		}
	})
}
