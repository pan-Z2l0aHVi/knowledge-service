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
	total, err := dao.FindFeedCount(ctx)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	feedList, err := dao.FindFeedList(ctx, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	list := make([]FeedListItem, len(feedList))
	for i, item := range feedList {
		author, err := dao.FindByAuthorID(ctx, item.AuthorID)
		if err != nil {
			tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
			return
		}
		list[i] = FeedListItem{
			Feed:       item,
			AuthorInfo: author,
		}
	}
	res := GetFeedListResp{
		Total: total,
		List:  list,
	}
	tools.RespSuccess(ctx, res)
}
