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

func (e *MaterialDAO) findMaterial(ctx *gin.Context, materialID string) (Material, error) {
	collection := e.GetDB().Collection("material")
	objID, err := primitive.ObjectIDFromHex(materialID)
	if err != nil {
		return Material{}, err
	}
	filter := bson.D{{"_id", objID}}
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
