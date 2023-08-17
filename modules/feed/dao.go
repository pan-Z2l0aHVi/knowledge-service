package feed

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedDao struct {
	*tools.Mongo
}

func (e *FeedDao) FindFeedList(ctx *gin.Context, page int, pageSize int) ([]Feed, error) {
	collection := e.GetDB().Collection("doc")
	filter := bson.D{{Key: "public", Value: true}}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}
	var feedList []Feed
	if err := cursor.All(ctx, &feedList); err != nil {
		return nil, err
	}
	if feedList == nil {
		feedList = []Feed{}
	}
	return feedList, nil
}

func (e *FeedDao) FindFeedCount(ctx *gin.Context) (int64, error) {
	collection := e.GetDB().Collection("doc")
	filter := bson.D{{Key: "public", Value: true}}
	return collection.CountDocuments(ctx, filter)
}
