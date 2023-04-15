package doc

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Doc) getInfo(ctx *gin.Context) {
	var params GetInfoParams
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := DocDAO{}
	docInfo, err := dao.findDoc(ctx, params.DocID)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	res := GetInfoResp{
		DocID: docInfo.ID,
		Doc: Doc{
			Content: docInfo.Content,
			Author:  docInfo.Author,
		},
	}
	tools.RespSuccess(ctx, res)
}

func (e *Doc) create(ctx *gin.Context) {

	// now := time.Now()
	// newUser := bson.D{
	// 	{Key: "nickname", Value: "122"},
	// 	{Key: "created_at", Value: now},
	// 	{Key: "updated_at", Value: now},
	// }
	// res, err := collection.InsertOne(context.TODO(), newUser)
}

func (e *Doc) update(ctx *gin.Context) {

}
