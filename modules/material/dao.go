package material

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MaterialDAO struct {
	*tools.Mongo
}

func (e *MaterialDAO) Find(ctx *gin.Context, materialID string) (Material, error) {
	collection := e.GetDB().Collection("material")
	objID, err := primitive.ObjectIDFromHex(materialID)
	if err != nil {
		return Material{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return Material{}, err
	}

	var materialInfo Material
	if err := res.Decode(&materialInfo); err != nil {
		return Material{}, err
	}
	return materialInfo, nil
}

func (e *MaterialDAO) Search(
	ctx *gin.Context,
	material_type int,
	keywords string,
	page int,
	pageSize int,
) ([]Material, error) {
	collection := e.GetDB().Collection("material")
	filter := bson.D{
		{Key: "type", Value: material_type},
		{Key: "name", Value: bson.D{
			{Key: "$regex", Value: keywords},
			{Key: "$options", Value: "i"},
		}},
	}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}
	var materialList []Material
	if err := cursor.All(ctx, &materialList); err != nil {
		return nil, err
	}
	if materialList == nil {
		materialList = []Material{}
	}
	return materialList, nil
}

func (e *MaterialDAO) SearchCount(ctx *gin.Context, material_type int, keywords string) (int64, error) {
	collection := e.GetDB().Collection("material")
	filter := bson.D{
		{Key: "type", Value: material_type},
		{Key: "name", Value: bson.D{
			{Key: "$regex", Value: keywords},
			{Key: "$options", Value: "i"},
		}},
	}
	return collection.CountDocuments(ctx, filter)
}

func (e *MaterialDAO) Create(ctx *gin.Context, material_type int, url string) (Material, error) {
	collection := e.GetDB().Collection("material")
	material := Material{
		ID:         primitive.NewObjectID(),
		URL:        url,
		Type:       material_type,
		UploaderID: "",
	}
	collection.InsertOne(ctx, material)
	return material, nil
}
