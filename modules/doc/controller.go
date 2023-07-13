package doc

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Doc) GetInfo(ctx *gin.Context) {
	var query GetInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := DocDAO{}
	docInfo, err := dao.Find(ctx, query.DocID)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	res := GetInfoResp{
		Doc: docInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Doc) Create(ctx *gin.Context) {
	dao := DocDAO{}
	docInfo, err := dao.Create(ctx)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, docInfo)
}

func (e *Doc) Update(ctx *gin.Context) {

}
