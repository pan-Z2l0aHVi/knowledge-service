package doc

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocDAO struct {
	*tools.Mongo
}

func (e *DocDAO) findDoc(ctx *gin.Context, docID string) (Doc, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return Doc{}, err
	}
	filter := bson.D{{"_id", objID}}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return Doc{}, err
	}

	var docInfo Doc
	if err := res.Decode(&docInfo); err != nil {
		return Doc{}, err
	}
	return docInfo, nil
}
