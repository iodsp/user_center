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
	"strings"
	"time"
)

func Store(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/store", func(c *gin.Context) {
		var param params.ResourceParams
		err := c.BindJSON(&param)
		param.Url = strings.Trim(param.Url, " ")
		name := strings.Trim(param.Name, " ")
		param.Url = strings.Trim(param.Url, " ")
		param.Desc = strings.Trim(param.Desc, " ")
		resource := service.NewResource(conf)
		domainService := service.NewDomain(conf)
		myLogger := my_log.NewLog(conf).Logger

		data, err1 := json.Marshal(param)
		if err1 != nil {
			myLogger.Info("Json marshaling failed：%s", err1)
		}
		myLogger.Info("params: " + string(data))

		if err == nil {
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

			resourceByUrl := resource.ShowByUrl(param.Url)
			if resourceByUrl.Id != 0 {
				myLogger.Info(common.DuplicatedUrlMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DuplicatedUrlMsg)
				return
			}

			resourceByName := resource.ShowByName(param.Name)
			if resourceByName.Id != 0 {
				myLogger.Info(common.NameUniqueMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
				return
			}

			if resourceByUrl.Id != 0 {
				myLogger.Info(common.DuplicatedUrlMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DuplicatedUrlMsg)
				return
			}

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

		myLogger.Info("Resource info id: " + stringId)

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
		domainService := service.NewDomain(conf)

		data, err1 := json.Marshal(param)
		if err1 != nil {
			myLogger.Info("Json marshaling failed：%s", err1)
		}
		myLogger.Info("params: " + string(data))

		myLogger.Info("User update id: " + idString)

		if err == nil {
			if param.DomainId == 0 && param.Name == "" && param.Url == "" && param.Desc == "" {
				myLogger.Info(common.NothingToUpdate)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NothingToUpdate)
				return
			}

			resourceInfo := resource.Show(id)
			//updating Resource does not exit
			if 0 == resourceInfo.Id {
				myLogger.Info(common.RecordNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RecordNotFoundMsg)
				return
			}

			name := strings.Trim(param.Name, " ")
			if name != "" {
				resourceByNameNotId := resource.ShowResourceByNameNotId(name, id)
				if resourceByNameNotId.Id != 0 {
					myLogger.Info(common.NameUniqueMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
					return

				}
				resourceInfo.Name = param.Name
			}

			if param.DomainId != 0 {
				domainInfo := domainService.ShowDomain(param.DomainId)
				if domainInfo.Id == 0 {
					myLogger.Info(common.DomainNotFoundMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
					return
				}
				resourceInfo.DomainId = param.DomainId
				resourceInfo.DomainName = domainInfo.Name
			}

			param.Url = strings.Trim(param.Url, " ")
			if param.Url != "" {
				resourceByUrlNotId := resource.ShowResourceByUrlNotId(param.Url, id)
				if resourceByUrlNotId.Id != 0 {
					myLogger.Info(common.DuplicatedUrlMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DuplicatedUrlMsg)
					return
				}
				resourceInfo.Url = param.Url
			}

			param.Desc = strings.Trim(param.Desc, " ")
			if param.Desc != "" {
				resourceInfo.Desc = param.Desc
			}

			resourceInfo.UpdatedAt = time.Now()

			data1, err2 := json.Marshal(resourceInfo)
			if err2 != nil {
				myLogger.Info("Json marshaling failed：%s", err2)
			}
			myLogger.Info("update params: " + string(data1))

			updateErr := resource.Update(resourceInfo)
			if updateErr == nil {
				myLogger.Info(common.UpdateSuccessMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.UpdateSuccessMsg)
			} else {
				myLogger.Error(common.UpdateSuccessMsg)
				apis.FormatResponseWithoutData(c, common.FailureCode, common.UpdateSuccessMsg)
			}
		} else {
			myLogger.Error(common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func DeleteResource(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/delete/:id", func(c *gin.Context) {
		Resource := service.NewResource(conf)
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		myLogger := my_log.NewLog(conf).Logger

		myLogger.Info("Delete resource id: " + stringId)

		if id == 0 {
			myLogger.Info(common.RecordNotFoundMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RecordNotFoundMsg)
			return
		}

		ResourceInfo := Resource.Show(id)

		//record not found
		if 0 == ResourceInfo.Id {
			myLogger.Info(common.RecordNotFoundMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RecordNotFoundMsg)
			return
		}

		delError := Resource.Delete(ResourceInfo)
		if delError == nil {
			myLogger.Info(common.DeleteSuccessMsg)
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
		} else {
			myLogger.Error(common.DeleteFailureMsg)
			apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
		}
	})
}
