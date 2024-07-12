package controller

import (
	"fmt"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type CommonController struct{}

// @Summary 埋点上报
// @Description 使用但不校验登录态
// @Produce json
// @Param request body entity.ReportPayload true "token和埋点数据"
// @Success 204 "ok"
// @Router /common/report [post]
func (e *CommonController) Report(ctx *gin.Context) {
	var payload entity.ReportPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	for i := range payload.Data {
		obj, ok := payload.Data[i].(map[string]interface{})
		if !ok {
			continue
		}
		obj["uid"] = userID
		obj["ip"] = ctx.ClientIP()
		obj["ua"] = ctx.Request.UserAgent()
		obj["date"] = time.UnixMilli(int64(obj["timestamp"].(float64)))
	}
	commonD := dao.CommonDAO{}
	if err := commonD.InsertReport(ctx, payload.Data); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary 获取统计数据
// @Description
// @Produce json
// @Param query query entity.GetStaticsQuery false "起止时间"
// @Success 200 {object} entity.GetStaticsResp "ok" "统计数据"
// @Router /common/statics [get]
func (e *CommonController) GetStatics(ctx *gin.Context) {
	var query entity.GetStaticsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	commonD := dao.CommonDAO{}
	pvCount, err := commonD.FindPVCount(ctx, query.StartTimestamp, query.EndTimestamp)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	uvCount, err := commonD.FindUVCount(ctx, query.StartTimestamp, query.EndTimestamp)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.GetStaticsResp{
		PV: pvCount,
		UV: uvCount,
	}
	tools.RespSuccess(ctx, res)
}

// @Summary 获取七牛云token
// @Description
// @Produce json
// @Success 200 {object} entity.GetBucketTokenResp "ok" "七牛云token"
// @Router /common/qiniu_token [get]
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

// @Summary 获取Cloudflare R2预上传链接
// @Description
// @Produce json
// @Success 200 {object} entity.GetSignedURLResp "ok" "Cloudflare R2预上传链接"
// @Router /common/r2_signed_url [get]
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
