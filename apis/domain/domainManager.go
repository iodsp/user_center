package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/apis"
	"github.com/iodsp/user_center/common"
	"github.com/iodsp/user_center/context"
	"github.com/iodsp/user_center/params"
	"github.com/iodsp/user_center/user_center"
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
		domainService := user_center.NewDomain(conf.Db(), conf.Debug())

		if err == nil {
			if name == "" {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameEmptyMsg)
				return
			}

			if domainType != 1 && domainType != 2 {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainTypeErrorMsg)
				return
			}

			//check duplicate name
			domainInfo := domainService.ShowDomainByName(name)
			if 0 != domainInfo.Id {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
				return
			}

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
	router.GET("/show/:id", func(c *gin.Context) {
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		domainService := user_center.NewDomain(conf.Db(), conf.Debug())

		//record not found
		domainInfo := domainService.ShowDomain(id)
		if 0 == domainInfo.Id{
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
			return
		}

		domain := domainService.ShowDomain(id)
		apis.FormatResponse(c, common.SuccessCode, "", &params.DomainItem{
			Id: domain.Id,
			Name: domain.Name,
			Type: domain.Type,
			UpdatedAt: domain.UpdatedAt,
			CreatedAt: domain.CreatedAt,
		})
	})
}

func List(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/list", func(c *gin.Context) {
		domainService := user_center.NewDomain(conf.Db(), conf.Debug())
		domain := domainService.DomainList()
		apis.FormatResponse(c, common.SuccessCode, "", domain)
	})
}

func Update(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/update/:id", func(c *gin.Context) {
		var param params.UpdateDomainParams
		err := c.BindJSON(&param);
		idString := c.Param("id")
		id, _ := strconv.Atoi(idString)
		name := param.Name
		domainType := param.Type
		domainService := user_center.NewDomain(conf.Db(), conf.Debug())

		if err == nil {
			domain := domainService.ShowDomain(id)
			//record not found
			if 0 == domain.Id {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
				return
			}

			if name != "" {
				//check duplicate name
				domainInfo := domainService.ShowDomainByNameNotId(name, id)
				if 0 != domainInfo.Id {
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
					return
				}
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
				//c.JSON(200, gin.H{"er": upError})
				apis.FormatResponseWithoutData(c, common.FailureCode, common.UpdateFailureMsg)
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
	router.POST("/delete/:id", func(c *gin.Context) {
		idString := c.Param("id")
		id, _ := strconv.Atoi(idString)
		domainService := user_center.NewDomain(conf.Db(), conf.Debug())

		//record not found
		domain := domainService.ShowDomain(id)
		if 0 == domain.Id {
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.DomainNotFoundMsg)
			return
		}

		delError := domainService.DeleteDomain(domain)
		if delError == nil {
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
		} else {
			apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
		}
	})
}
