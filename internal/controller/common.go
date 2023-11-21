package controller

import (
	"fmt"
	"knowledge-service/internal/api"
	"knowledge-service/internal/dao"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	QINIU_ACCESS_KEY = "ZJAw5p66HbXlJQbXjDV5Y_qLAQXEirlm8MXcG-l2"
	QINIU_SECRET_KEY = "JOR_yrHI5nxg1SbZ1DF0i0BmkHYW45UH9FAXXx3m"
	QINIU_BUCKET     = "youknown1120"
)
const (
	R2_ACCESS_ID     = "eff6f464c82bf6f773a57d4e5428ad4e"
	R2_ACCESS_SECRET = "44f4bee6c901a4c7793f462a1f9941091101f1bf11b50778d1a22a0e29865608"
	R2_ACCOUNT_ID    = "70bc20cd210d1c9e762acb3786056b90"
	R2_BUCKET        = "youknown"
)

type CommonController struct{}

func (e *CommonController) Report(ctx *gin.Context) {
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
	commonD := dao.CommonDao{}
	if err := commonD.InsertReport(ctx, jsonData); err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (e *CommonController) GetQiniuToken(ctx *gin.Context) {
	putPolicy := storage.PutPolicy{
		Scope: QINIU_BUCKET,
	}
	putPolicy.Expires = 86400
	mac := qbox.NewMac(QINIU_ACCESS_KEY, QINIU_SECRET_KEY)
	upToken := putPolicy.UploadToken(mac)
	res := api.GetBucketTokenResp{
		Token: upToken,
	}
	tools.RespSuccess(ctx, res)
}

func (e *CommonController) GetR2SignedURL(ctx *gin.Context) {
	var query api.GetSignedURLQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", R2_ACCOUNT_ID),
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(R2_ACCESS_ID, R2_ACCESS_SECRET, ""),
		),
	)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)
	presignResult, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(R2_BUCKET),
		Key:    aws.String(query.Key),
	})
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := api.GetSignedURLResp{
		URL: presignResult.URL,
	}
	tools.RespSuccess(ctx, res)
}
