package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS int = 0
	TOAST   int = 1
)

func RespSuccess(ctx *gin.Context, val interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": SUCCESS,
		"msg":  "成功",
		"data": val,
	})
}

func RespFail(ctx *gin.Context, code int, msg string, val interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": val,
	})
}
