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

const (
	CollectionName = "doc"
)

func (e *DocDAO) Find(ctx *gin.Context, docID string) (Doc, error) {
	collection := e.GetDB().Collection(CollectionName)
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return Doc{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}
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

func (e *DocDAO) Create(ctx *gin.Context) (Doc, error) {
	collection := e.GetDB().Collection(CollectionName)
	doc := Doc{
		ID:       primitive.NewObjectID(),
		AuthorID: "qwer",
		Content:  "<p>content</p>",
	}
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return Doc{}, err
	}
	return doc, nil
}
