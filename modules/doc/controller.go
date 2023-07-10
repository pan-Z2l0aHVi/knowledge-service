package doc

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Doc) getInfo(ctx *gin.Context) {
	var params GetInfoQuery
	if err := ctx.ShouldBindQuery(&params); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := DocDAO{}
	docInfo, err := dao.find(ctx, params.DocID)
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
	}
	tools.RespSuccess(ctx, docInfo)
}

func (e *Doc) update(ctx *gin.Context) {

}
