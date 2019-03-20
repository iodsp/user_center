package domain

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

//r.POST("/domain/store", domain.Store)
func Store(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/store", func(c *gin.Context) {
		var param params.DomainParams
		err := c.BindJSON(&param)
		name := param.Name
		domainType := param.Type
		domainService := service.NewDomain(conf)
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

			if domainType != 1 && domainType != 2 {
				myLogger.Info(common.DomainTypeErrorMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainTypeErrorMsg)
				return
			}

			//check duplicate name
			domainInfo := domainService.ShowDomainByName(name)
			if 0 != domainInfo.Id {
				myLogger.Info(common.NameUniqueMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
				return
			}

			insertErr := domainService.StoreDomain(param)
			if insertErr == nil {
				myLogger.Info(common.StoreDomainSuccess)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.StoreDomainSuccess)
			} else {
				myLogger.Error(common.StoreDomainFailure)
				apis.FormatResponseWithoutData(c, common.FailureCode, common.StoreDomainFailure)
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
		domainService := service.NewDomain(conf)
		myLogger := my_log.NewLog(conf).Logger
		myLogger.Info("Domain info id: " + stringId)

		//record not found
		domainInfo := domainService.ShowDomain(id)
		if 0 == domainInfo.Id {
			myLogger.Info(common.DomainNotFoundMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
			return
		}

		domain := domainService.ShowDomain(id)
		apis.FormatResponse(c, common.SuccessCode, "", &params.DomainItem{
			Id:        domain.Id,
			Name:      domain.Name,
			Type:      domain.Type,
			UpdatedAt: domain.UpdatedAt,
			CreatedAt: domain.CreatedAt,
		})
	})
}

func List(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/list", func(c *gin.Context) {
		domainService := service.NewDomain(conf)
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
		domainService := service.NewDomain(conf)
		myLogger := my_log.NewLog(conf).Logger

		data, err1 := json.Marshal(param)
		if err1 != nil {
			myLogger.Info("Json marshaling failed：%s", err1)
		}
		myLogger.Info("params " + string(data))

		if err == nil {
			domain := domainService.ShowDomain(id)
			//record not found
			if 0 == domain.Id {
				myLogger.Info(common.DomainNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
				return
			}

			if name != "" {
				//check duplicate name
				domainInfo := domainService.ShowDomainByNameNotId(name, id)
				if 0 != domainInfo.Id {
					myLogger.Info(common.NameUniqueMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
					return
				}
				domain.Name = name
			} else if domainType == 0 && name == "" {
				myLogger.Info(common.NothingToUpdate)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NothingToUpdate)
				return
			}

			if domainType != 0 {
				domain.Type = domainType
			}

			if domainType != 1 && domainType != 2 {
				myLogger.Info(common.DomainTypeErrorMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainTypeErrorMsg)
				return
			}

			domain.UpdatedAt = time.Now()
			upError := domainService.UpdateDomain(domain)
			if upError != nil {
				myLogger.Info(common.UpdateFailureMsg)
				apis.FormatResponseWithoutData(c, common.FailureCode, common.UpdateFailureMsg)
				return
			} else {
				myLogger.Error(common.UpdateSuccessMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.UpdateSuccessMsg)
				return
			}
		} else {
			myLogger.Error(common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func DeleteDomain(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/delete/:id", func(c *gin.Context) {
		idString := c.Param("id")
		id, _ := strconv.Atoi(idString)
		domainService := service.NewDomain(conf)
		myLogger := my_log.NewLog(conf).Logger

		myLogger.Info("Delete domain id: " + idString)

		//record not found
		domain := domainService.ShowDomain(id)
		if 0 == domain.Id {
			myLogger.Info(common.DomainNotFoundMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.DomainNotFoundMsg)
			return
		}

		delError := domainService.DeleteDomain(domain)
		if delError == nil {
			myLogger.Info(common.DeleteSuccessMsg)
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
		} else {
			myLogger.Info(common.DeleteFailureMsg)
			apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
		}
	})
}
