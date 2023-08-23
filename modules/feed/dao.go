package feed

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	sort := bson.D{{Key: "update_time", Value: -1}}
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  sort,
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

func (e *FeedDao) FindByAuthorID(ctx *gin.Context, userID string) (Author, error) {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return Author{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return Author{}, err
	}
	var author Author
	if err := res.Decode(&author); err != nil {
		return Author{}, err
	}
	return author, nil
}
