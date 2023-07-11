package doc

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Doc) getInfo(ctx *gin.Context) {
	var query GetInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := DocDAO{}
	docInfo, err := dao.find(ctx, query.DocID)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	res := GetInfoResp{
		Doc: docInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Doc) create(ctx *gin.Context) {
	dao := DocDAO{}
	docInfo, err := dao.create(ctx)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, docInfo)
}

func (e *Doc) update(ctx *gin.Context) {

}
