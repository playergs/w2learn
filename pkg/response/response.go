package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	SuccessCodeDefault = 0
	ErrorCodeDefault   = -1
)

const (
	SuccessMsgDefault = "success"
	ErrorMsgDefault   = "error"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SuccessCodeDefault,
		Msg:  SuccessMsgDefault,
		Data: data,
	})
}

func Error(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: ErrorCodeDefault,
		Msg:  ErrorMsgDefault,
		Data: data,
	})
}
