package common

import (
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func (e *Common) Report(ctx *gin.Context) {
	var jsonData []interface{}
	if err := ctx.ShouldBindJSON(&jsonData); err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	for i := range jsonData {
		obj, ok := jsonData[i].(map[string]interface{})
		if !ok {
			continue
		}
		obj["ip"] = ctx.ClientIP()
		obj["ua"] = ctx.Request.UserAgent()
	}
	dao := CommonDao{}
	dao.insertReport(ctx, jsonData)
	ctx.Status(http.StatusNoContent)
}

func (e *Common) GetQiniuToken(ctx *gin.Context) {
	accessKey := "ZJAw5p66HbXlJQbXjDV5Y_qLAQXEirlm8MXcG-l2"
	secretKey := "JOR_yrHI5nxg1SbZ1DF0i0BmkHYW45UH9FAXXx3m"
	bucket := "youknown1120"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	putPolicy.Expires = 86400
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	res := GetBucketTokenResp{
		Token: upToken,
	}
	tools.RespSuccess(ctx, res)
}
