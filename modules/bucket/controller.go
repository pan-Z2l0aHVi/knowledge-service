package bucket

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func (e *Bucket) GetQiniuToken(ctx *gin.Context) {
	accessKey := "ZJAw5p66HbXlJQbXjDV5Y_qLAQXEirlm8MXcG-l2"
	secretKey := "JOR_yrHI5nxg1SbZ1DF0i0BmkHYW45UH9FAXXx3m"
	bucket := "newyouknown"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	putPolicy.Expires = 86400
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	res := GetTokenResp{
		Token: upToken,
	}
	tools.RespSuccess(ctx, res)
}
