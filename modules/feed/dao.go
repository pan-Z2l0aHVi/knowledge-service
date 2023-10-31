package feed

import (
	"knowledge-base-service/tools"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedDao struct {
	*tools.Mongo
}

func (e *FeedDao) FindFeedList(
	ctx *gin.Context,
	page int,
	pageSize int,
	keywords string,
	sortBy string,
	asc int,
	authorID string,
) ([]Feed, error) {
	collection := e.GetDB().Collection("doc")
	filter := bson.M{"public": true}
	if authorID != "" {
		filter["author_id"] = authorID
	}
	if keywords != "" {
		keywords := regexp.QuoteMeta(keywords)
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": keywords, "$options": "i"}},
			{"summary": bson.M{"$regex": keywords, "$options": "i"}},
		}
	}
	sort := bson.M{}
	if sortBy != "" && asc != 0 {
		sort[sortBy] = asc
	}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  sort,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var feedList []Feed
	if err := cursor.All(ctx, &feedList); err != nil {
		return nil, err
	}
	if feedList == nil {
		feedList = []Feed{}
	}

	clone := make([]Feed, len(feedList))
	for i, feed := range feedList {
		if feed.Likes == nil {
			feed.Likes = []Like{}
		}
		feed.LikesCount = len(feed.Likes)
		clone[i] = feed
	}
	return clone, nil
}

func (e *FeedDao) FindFeed(ctx *gin.Context, feedID string) (Feed, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return Feed{}, err
	}
	filter := bson.M{"_id": objID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return Feed{}, err
	}
	var feed Feed
	if err := res.Decode(&feed); err != nil {
		return Feed{}, err
	}
	return feed, nil
}

func (e *FeedDao) FindAuthorByID(ctx *gin.Context, userID string) (Author, error) {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return Author{}, err
	}
	filter := bson.M{"_id": objID}
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

func (e *FeedDao) CheckLiked(ctx *gin.Context, userID string, feedID string) (bool, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": objID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return false, err
	}
	var feed Feed
	if err := res.Decode(&feed); err != nil {
		return false, err
	}
	likes := feed.Likes
	for _, like := range likes {
		if like.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (e *FeedDao) Like(ctx *gin.Context, userID string, feedID string) error {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{
			"likes": bson.M{
				"user_id":       userID,
				"creation_time": time.Now(),
			},
		},
	}
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (e *FeedDao) UnLike(ctx *gin.Context, userID string, feedID string) error {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$pull": bson.M{
			"likes": bson.M{"user_id": userID},
		},
	}
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}
