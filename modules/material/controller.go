package material

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Material) getInfo(ctx *gin.Context) {
	var params GetInfoParams
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := MaterialDAO{}
	materialInfo, err := dao.findMaterial(ctx, params.MaterialID)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	res := GetInfoResp{
		MaterialID: materialInfo.ID,
		Material: Material{
			URL:      materialInfo.URL,
			Uploader: materialInfo.Uploader,
		},
	}
	tools.RespSuccess(ctx, res)
}

func (e *Material) create(ctx *gin.Context) {

	// now := time.Now()
	// newUser := bson.D{
	// 	{Key: "nickname", Value: "122"},
	// 	{Key: "created_at", Value: now},
	// 	{Key: "updated_at", Value: now},
	// }
	// res, err := collection.InsertOne(context.TODO(), newUser)
}

func (e *Material) update(ctx *gin.Context) {

}

func (e *Material) upload(ctx *gin.Context) {

}
