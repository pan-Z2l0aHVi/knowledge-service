package dao

import (
	"knowledge-service/internal/model"
	"knowledge-service/pkg/tools"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MaterialDAO struct {
	*tools.Mongo
}

func (e *MaterialDAO) Find(ctx *gin.Context, materialID string) (model.Material, error) {
	collection := e.GetDB().Collection("material")
	objID, err := primitive.ObjectIDFromHex(materialID)
	if err != nil {
		return model.Material{}, err
	}
	filter := bson.M{"_id": objID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return model.Material{}, err
	}

	var materialInfo model.Material
	if err := res.Decode(&materialInfo); err != nil {
		return model.Material{}, err
	}
	return materialInfo, nil
}

func (e *MaterialDAO) FindList(
	ctx *gin.Context,
	materialType int,
	keywords string,
	page int,
	pageSize int,
) ([]model.Material, error) {
	collection := e.GetDB().Collection("material")
	filter := bson.M{
		"type": materialType,
		"name": bson.M{
			"$regex":   regexp.QuoteMeta(keywords),
			"$options": "i",
		},
	}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	sort := bson.M{"update_time": -1}
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  sort,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var materialList []model.Material
	if err := cursor.All(ctx, &materialList); err != nil {
		return nil, err
	}
	if materialList == nil {
		materialList = []model.Material{}
	}
	return materialList, nil
}

func (e *MaterialDAO) Create(ctx *gin.Context, materialType int, url string) (model.Material, error) {
	collection := e.GetDB().Collection("material")
	material := model.Material{
		ID:         primitive.NewObjectID(),
		URL:        url,
		Type:       materialType,
		UploaderID: "",
	}
	_, err := collection.InsertOne(ctx, material)
	if err != nil {
		return model.Material{}, err
	}
	return material, nil
}
