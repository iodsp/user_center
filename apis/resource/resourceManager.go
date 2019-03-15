package resource

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
		var param params.ResourceParams
		err := c.BindJSON(&param)
		name := param.Name
		resource := service.NewResource(conf)
		domainService := service.NewDomain(conf)
		myLogger := my_log.NewLog(conf).Logger
		data, err := json.Marshal(param)
		if err != nil {
			myLogger.Info("Json marshaling failed：%s", err)
		}
		myLogger.Info("params: " + string(data))

		if name == "" {
			myLogger.Info(common.NameEmptyMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameEmptyMsg)
			return
		}

		if param.DomainId == 0 {
			myLogger.Info(common.DomainIdEmptyMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainIdEmptyMsg)
			return
		}

		if param.Url == "" {
			myLogger.Info(common.UrlEmptyMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UrlEmptyMsg)
			return
		}

		if err == nil {
			domainInfo := domainService.ShowDomain(param.DomainId)
			if domainInfo.Id == 0 {
				myLogger.Info(common.DomainNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
				return
			}
			param.DomainName = domainInfo.Name

			//insert a new record
			insertErr := resource.StoreResource(param)

			if insertErr == nil {
				myLogger.Info(common.SaveSuccessMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.SaveSuccessMsg)
			} else {
				myLogger.Error(insertErr)
				apis.FormatResponseWithoutData(c, common.FailureCode, common.SaveFailureMsg)
			}
		} else {
			myLogger.Error(err)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func Show(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/show/:id", func(c *gin.Context) {
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		resource := service.NewResource(conf)
		myLogger := my_log.NewLog(conf).Logger

		resourceInfo := resource.Show(id)

		//record does not exist
		if 0 == resourceInfo.Id {
			myLogger.Info(common.RecordNotFoundMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RecordNotFoundMsg)
			return
		}

		apis.FormatResponse(c, common.SuccessCode, "", &params.ResourceItem{
			Id:         resourceInfo.Id,
			DomainId:   resourceInfo.DomainId,
			DomainName: resourceInfo.DomainName,
			Name:       resourceInfo.Name,
			Url:        resourceInfo.Url,
			Desc:       resourceInfo.Desc,
			CreatedAt:  resourceInfo.CreatedAt,
			UpdatedAt:  resourceInfo.UpdatedAt,
		})
	})
}

func List(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/list", func(c *gin.Context) {
		resource := service.NewResource(conf)
		result := resource.List()
		apis.FormatResponse(c, common.SuccessCode, "", result)
	})
}

func Update(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/update/:id", func(c *gin.Context) {
		var param params.ResourceUpdateParams
		idString := c.Param("id")
		id, _ := strconv.Atoi(idString)
		err := c.BindJSON(&param)
		resource := service.NewResource(conf)
		myLogger := my_log.NewLog(conf).Logger

		if param.DomainId == 0 && param.Name == "" && param.Url == "" && param.Desc == "" {
			myLogger.Info(common.NothingToUpdate)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NothingToUpdate)
			return
		}

		if err == nil {
			ResourceInfo := resource.Show(id)
			//updating Resource does not exit
			if 0 == ResourceInfo.Id {
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RecordNotFoundMsg)
				return
			}
			name := param.Name

			if name == "" {
				apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.NothingToUpdate)
				return
			} else {
				ResourceInfo.Name = param.Name
				ResourceInfo.UpdatedAt = time.Now()
			}

			updateErr := resource.Update(ResourceInfo)
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

func DeleteResource(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/delete/:id", func(c *gin.Context) {
		Resource := service.NewResource(conf)
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		ResourceInfo := Resource.Show(id)

		//record not found
		if 0 == ResourceInfo.Id {
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RecordNotFoundMsg)
			return
		}

		delError := Resource.Delete(ResourceInfo)
		if delError == nil {
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
		} else {
			apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
		}
	})
}
