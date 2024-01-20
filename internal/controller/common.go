package controller

import (
	"fmt"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
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

type CommonController struct{}

func (e *CommonController) Report(ctx *gin.Context) {
	var jsonData []interface{}
	if err := ctx.ShouldBindJSON(&jsonData); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
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
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (e *CommonController) GetQiniuToken(ctx *gin.Context) {
	putPolicy := storage.PutPolicy{
		Scope: consts.QINIU_BUCKET,
	}
	putPolicy.Expires = 86400
	mac := qbox.NewMac(consts.QINIU_ACCESS_KEY, consts.QINIU_SECRET_KEY)
	upToken := putPolicy.UploadToken(mac)
	res := entity.GetBucketTokenResp{
		Token: upToken,
	}
	tools.RespSuccess(ctx, res)
}

func (e *CommonController) GetR2SignedURL(ctx *gin.Context) {
	var query entity.GetSignedURLQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", consts.R2_ACCOUNT_ID),
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(consts.R2_ACCESS_ID, consts.R2_ACCESS_SECRET, ""),
		),
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)
	presignResult, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(consts.R2_BUCKET),
		Key:    aws.String(query.Key),
	})
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.GetSignedURLResp{
		URL: presignResult.URL,
	}
	tools.RespSuccess(ctx, res)
}
