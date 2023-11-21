package tools

import (
	"knowledge-service/pkg/consts"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespSuccess(ctx *gin.Context, val interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": consts.SuccessCode,
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
