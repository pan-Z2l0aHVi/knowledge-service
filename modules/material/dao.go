package material

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaterialDAO struct {
	*tools.Mongo
}

const (
	COLLECTION_NAME = "material"
)

func (e *MaterialDAO) find(ctx *gin.Context, materialID string) (Material, error) {
	collection := e.GetDB().Collection(COLLECTION_NAME)
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

func (e *MaterialDAO) search(ctx *gin.Context, material_type int, keywords string) ([]Material, error) {
	collection := e.GetDB().Collection(COLLECTION_NAME)
	filter := bson.D{
		{Key: "type", Value: material_type},
		{Key: "name", Value: bson.D{
			{Key: "$regex", Value: keywords},
			{Key: "$options", Value: "i"},
		}},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var materialList []Material
	if err := cursor.All(ctx, &materialList); err != nil {
		return nil, err
	}
	return materialList, nil
}

func (e *MaterialDAO) create(ctx *gin.Context) (Material, error) {
	collection := e.GetDB().Collection(COLLECTION_NAME)
	material := Material{
		ID:         primitive.NewObjectID(),
		URL:        "",
		Type:       1,
		UploaderID: "",
	}
	collection.InsertOne(ctx, material)
	return material, nil
}
