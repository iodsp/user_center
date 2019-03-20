package user

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/apis"
	"github.com/iodsp/user_center/common"
	"github.com/iodsp/user_center/context"
	"github.com/iodsp/user_center/my_log"
	"github.com/iodsp/user_center/params"
	"github.com/iodsp/user_center/service"
	"regexp"
	"strconv"
	"time"
)

var coder = base64.NewEncoding(common.Base64Table)

func Store(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/store", func(c *gin.Context) {
		var param params.UserParams
		err := c.BindJSON(&param)
		userService := service.NewUser(conf)
		domainService := service.NewDomain(conf)
		name := param.Name
		phone := param.Phone
		userType := param.Type
		password := param.Password
		myLogger := my_log.NewLog(conf).Logger

		data, err := json.Marshal(param)
		if err != nil {
			myLogger.Info("Json marshaling failed：%s", err)
		}
		myLogger.Info("params: " + string(data))

		if err == nil {
			if name == "" {
				myLogger.Info(common.AddUserMsg + common.NameEmptyMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameEmptyMsg)
				return
			}
			//check duplicate name
			userInfo := userService.ShowByName(name)
			if 0 != userInfo.Id {
				myLogger.Info(common.AddUserMsg + common.NameUniqueMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
				return
			}

			if param.Phone == "" {
				myLogger.Info(common.AddUserMsg + common.PhoneEmptyMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.PhoneEmptyMsg)
				return
			}
			userInfoByPhone := userService.ShowByPhone(phone)
			if 0 != userInfoByPhone.Id {
				myLogger.Info(common.AddUserMsg + common.PhoneUniqueMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.PhoneUniqueMsg)
				return
			}
			if !isPhoneValidated(phone) {
				myLogger.Info(common.AddUserMsg + common.PhonePatternErrMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.PhonePatternErrMsg)
				return
			}

			if param.DomainId == 0 {
				myLogger.Info(common.AddUserMsg + common.DomainIdEmptyMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainIdEmptyMsg)
				return
			}

			if password == "" {
				myLogger.Info(common.AddUserMsg + common.PasswordEmptyMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.PasswordEmptyMsg)
				return
			}

			//base64 encode password
			param.Password = coder.EncodeToString([]byte(password))

			if param.RegisterPlatform == 0 {
				param.RegisterPlatform = common.DefaultRegisterPlatform
			} else if param.RegisterPlatform != 1 && param.RegisterPlatform != 2 {
				myLogger.Info(common.AddUserMsg + common.RegisterPlatformErrMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RegisterPlatformErrMsg)
				return
			}

			domainInfo := domainService.ShowDomain(param.DomainId)
			if 0 == domainInfo.Id {
				myLogger.Info(common.AddUserMsg + common.DomainNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
				return
			}
			param.DomainName = domainInfo.Name

			if userType == 0 {
				param.Type = common.DefaultUserType
			} else if userType != 1 && userType != 2 {
				myLogger.Info(common.AddUserMsg + common.UserTypeErrMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserTypeErrMsg)
				return
			}
			userInfo.UpdatedAt = time.Now()

			//log insert param
			data2, err2 := json.Marshal(param)
			if err2 != nil {
				myLogger.Info("Json marshaling failed：%s", err2)
			}
			myLogger.Info(common.AddUserMsg + "insert params" + string(data2))

			//insert a new record
			insertErr := userService.StoreUser(param)

			if insertErr == nil {
				myLogger.Info(common.AddUserMsg + common.SaveSuccessMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.SaveSuccessMsg)
			} else {
				myLogger.Error(common.AddUserMsg + insertErr.Error())
				apis.FormatResponseWithoutData(c, common.FailureCode, common.SaveFailureMsg)
			}
		} else {
			myLogger.Error(common.AddUserMsg + common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}

func isPhoneValidated(phone string) bool {
	regular := common.UserPhoneRegular
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

func Show(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/show/:id", func(c *gin.Context) {
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		userService := service.NewUser(conf)
		myLogger := my_log.NewLog(conf).Logger

		myLogger.Info("User info id: " + stringId)

		userInfo := userService.Show(id)

		//record does not exist
		if 0 == userInfo.Id {
			myLogger.Info(common.ShowUserMsg + common.UserNotFoundMsg)
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserNotFoundMsg)
			return
		}

		apis.FormatResponse(c, common.SuccessCode, "", &params.UserItem{
			Id:               userInfo.Id,
			Name:             userInfo.Name,
			Phone:            userInfo.Phone,
			RegisterPlatform: userInfo.RegisterPlatform,
			DomainName:       userInfo.DomainName,
			Type:             userInfo.Type,
			IdentityId:       userInfo.IdentityId,
			UpdatedAt:        userInfo.UpdatedAt,
			CreatedAt:        userInfo.CreatedAt,
		})
	})
}

func List(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/list", func(c *gin.Context) {
		userService := service.NewUser(conf)
		result := userService.List()
		apis.FormatResponse(c, common.SuccessCode, "", result)
	})
}

func DeleteUser(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/delete/:id", func(c *gin.Context) {
		user := service.NewUser(conf)
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		userInfo := user.Show(id)
		myLogger := my_log.NewLog(conf).Logger

		myLogger.Info("Delete user id: " + stringId)

		//record not found
		if 0 == userInfo.Id {
			apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserNotFoundMsg)
			return
		}

		delError := user.Delete(userInfo)
		if delError == nil {
			myLogger.Info(common.DeleteUserMsg + common.DeleteSuccessMsg)
			apis.FormatResponseWithoutData(c, common.SuccessCode, common.DeleteSuccessMsg)
		} else {
			myLogger.Error(common.DeleteUserMsg + common.DeleteFailureMsg)
			apis.FormatResponseWithoutData(c, common.FailureCode, common.DeleteFailureMsg)
		}
	})
}

func UpdateUser(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/update/:id", func(c *gin.Context) {
		var param params.UserUpdateParams
		stringId := c.Param("id")
		id, _ := strconv.Atoi(stringId)
		err := c.BindJSON(&param)
		userService := service.NewUser(conf)
		domainService := service.NewDomain(conf)
		name := param.Name
		phone := param.Phone
		userType := param.Type
		password := param.Password
		myLogger := my_log.NewLog(conf).Logger

		data, err := json.Marshal(param)
		if err != nil {
			myLogger.Info("Json marshaling failed：%s", err)
		}
		myLogger.Info(common.UpdateUserMsg + "params " + string(data))

		if err == nil {
			if id == 0 {
				myLogger.Info(common.UpdateUserMsg + common.IdEmpty)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.IdEmpty)
				return
			}
			myLogger.Info(common.UpdateUserMsg + "id: " + stringId)

			userInfo := userService.Show(id)
			logData, err := json.Marshal(userInfo)
			if err != nil {
				myLogger.Info("Json marshaling failed：%s", err)
			}
			myLogger.Info(common.UpdateUserMsg + "record to be updated: " + string(logData))

			if userInfo.Id == 0 {
				myLogger.Info(common.UpdateUserMsg + common.UserNotFoundMsg)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserNotFoundMsg)
				return
			}

			if name == "" && password == "" && phone == "" && param.RegisterPlatform == 0 && param.DomainId == 0 && userType == 0 {
				myLogger.Info(common.UpdateUserMsg + common.NothingToUpdate)
				apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NothingToUpdate)
				return
			}

			if name != "" {
				//check duplicate name
				userInfo1 := userService.UpdateShowByName(name, id)
				if 0 != userInfo1.Id {
					myLogger.Info(common.UpdateUserMsg + common.NameUniqueMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.NameUniqueMsg)
					return
				}
				userInfo.Name = name
			}

			//base64 encode password
			if password != "" {
				userInfo.Password = coder.EncodeToString([]byte(password))
			}

			if phone != "" {
				userInfoByPhone := userService.UpdateShowByPhone(phone, id)
				if 0 != userInfoByPhone.Id {
					myLogger.Info(common.UpdateUserMsg + common.PhoneUniqueMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.PhoneUniqueMsg)
					return
				}

				if !isPhoneValidated(phone) {
					myLogger.Info(common.UpdateUserMsg + common.PhonePatternErrMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.PhonePatternErrMsg)
					return
				}
				userInfo.Phone = phone
			}

			if param.RegisterPlatform != 0 {
				if param.RegisterPlatform != 1 && param.RegisterPlatform != 2 {
					myLogger.Info(common.UpdateUserMsg + common.RegisterPlatformErrMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.RegisterPlatformErrMsg)
					return
				}
				userInfo.RegisterPlatform = param.RegisterPlatform
			}

			if param.DomainId != 0 {
				domainInfo := domainService.ShowDomain(param.DomainId)
				if 0 == domainInfo.Id {
					myLogger.Info(common.UpdateUserMsg + common.DomainNotFoundMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.DomainNotFoundMsg)
					return
				}
				userInfo.DomainName = domainInfo.Name
				userInfo.DomainId = param.DomainId
			}

			if userType != 0 {
				if userType != 1 && userType != 2 {
					myLogger.Info(common.UpdateUserMsg + common.UserTypeErrMsg)
					apis.FormatResponseWithoutData(c, common.ParamErrorCode, common.UserTypeErrMsg)
					return
				}
				userInfo.Type = userType
			}

			//log update param
			data2, err := json.Marshal(userInfo)
			if err != nil {
				myLogger.Info("Json marshaling failed：%s", err)
			}
			myLogger.Info(common.UpdateUserMsg + "update params" + string(data2))

			//update a new record
			insertErr := userService.Update(userInfo)

			if insertErr == nil {
				myLogger.Info(common.UpdateUserMsg + common.UpdateSuccessMsg)
				apis.FormatResponseWithoutData(c, common.SuccessCode, common.UpdateSuccessMsg)
			} else {
				myLogger.Error(common.UpdateUserMsg + insertErr.Error())
				apis.FormatResponseWithoutData(c, common.FailureCode, common.UpdateFailureMsg)
			}
		} else {
			myLogger.Error(common.UpdateUserMsg + common.ParseParamErrorMsg)
			apis.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
		}
	})
}
