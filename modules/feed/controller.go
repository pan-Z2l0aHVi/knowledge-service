package feed

import (
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Feed) GetFeedList(ctx *gin.Context) {
	var query GetFeedListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := FeedDao{}
	feedList, err := dao.FindFeedList(ctx, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
	}
	total, err := dao.FindFeedCount(ctx)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
	}
	res := GetFeedListResp{
		Total: total,
		List:  feedList,
	}
	tools.RespSuccess(ctx, res)
}
