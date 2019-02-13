package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/common"
)

type result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type resultWithoutDate struct {
	Code    int  `json:"code"`
	Message string `json:"message"`
}

func FormatResponse(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(common.HTTP_SUCCESS_CODE, &result{
		code,
		msg,
		data,
	})
}

func FormatResponseWithoutData(c *gin.Context, code int, msg string){
	c.JSON(common.HTTP_SUCCESS_CODE, &resultWithoutDate{
		code,
		msg,
	})
	return
}
