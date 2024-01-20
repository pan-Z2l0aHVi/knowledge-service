package dao

import (
	"knowledge-service/internal/model"
	"knowledge-service/pkg/tools"
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

func (e *FeedDao) Create(ctx *gin.Context, creatorID string, subjectID string, subjectType string) (model.Feed, error) {
	collection := e.GetDB().Collection("feed")
	now := time.Now()
	feed := model.Feed{
		ID:           primitive.NewObjectID(),
		CreatorID:    creatorID,
		SubjectID:    subjectID,
		SubjectType:  subjectType,
		Likes:        []model.Like{},
		LikesCount:   0,
		CreationTime: now,
		UpdateTime:   now,
	}
	_, err := collection.InsertOne(ctx, feed)
	if err != nil {
		return model.Feed{}, err
	}
	return feed, nil
}

func (e *FeedDao) Update(ctx *gin.Context, feedID string) (model.Feed, error) {
	collection := e.GetDB().Collection("feed")
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return model.Feed{}, err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"update_time": time.Now()}}

	var feed model.Feed
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&feed); err != nil {
		return model.Feed{}, err
	}
	return feed, nil
}

func (e *FeedDao) DeleteMany(ctx *gin.Context, feedIDs []string) error {
	collection := e.GetDB().Collection("feed")
	var objIDs []primitive.ObjectID
	for _, feedID := range feedIDs {
		id, err := primitive.ObjectIDFromHex(feedID)
		if err != nil {
			return err
		}
		objIDs = append(objIDs, id)
	}
	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	if _, err := collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}

func (e *FeedDao) DeleteManyBySubject(ctx *gin.Context, subjectIDs []string, subjectType string) error {
	collection := e.GetDB().Collection("feed")
	filter := bson.M{
		"subject_id":   bson.M{"$in": subjectIDs},
		"subject_type": subjectType,
	}
	if _, err := collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}

func (e *FeedDao) FindList(
	ctx *gin.Context,
	page int,
	pageSize int,
	keywords string,
	sortBy string,
	asc int,
	authorID string,
) ([]model.Feed, error) {
	collection := e.GetDB().Collection("feed")
	filter := bson.M{}
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
	var feeds []model.Feed
	if err := cursor.All(ctx, &feeds); err != nil {
		return nil, err
	}
	if feeds == nil {
		feeds = []model.Feed{}
	}
	for i := range feeds {
		if feeds[i].Likes == nil {
			feeds[i].Likes = []model.Like{}
		}
		feeds[i].LikesCount = len(feeds[i].Likes)
	}
	return feeds, nil
}

func (e *FeedDao) Find(ctx *gin.Context, feedID string) (model.Feed, error) {
	collection := e.GetDB().Collection("feed")
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return model.Feed{}, err
	}
	filter := bson.M{"_id": objID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return model.Feed{}, err
	}
	var feed model.Feed
	if err := res.Decode(&feed); err != nil {
		return model.Feed{}, err
	}
	if feed.Likes == nil {
		feed.Likes = []model.Like{}
	}
	feed.LikesCount = len(feed.Likes)
	return feed, nil
}

func (e *FeedDao) FindBySubject(ctx *gin.Context, subjectID string, subjectType string) (model.Feed, error) {
	collection := e.GetDB().Collection("feed")
	filter := bson.M{
		"subject_id":   subjectID,
		"subject_type": subjectType,
	}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return model.Feed{}, err
	}
	var feed model.Feed
	if err := res.Decode(&feed); err != nil {
		return model.Feed{}, err
	}
	feed.LikesCount = len(feed.Likes)
	return feed, nil
}

func (e *FeedDao) CheckLiked(ctx *gin.Context, userID string, feedID string) (bool, error) {
	feed, err := e.Find(ctx, feedID)
	if err != nil {
		return false, err
	}
	for _, like := range feed.Likes {
		if like.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (e *FeedDao) Like(ctx *gin.Context, userID string, feedID string) error {
	collection := e.GetDB().Collection("feed")
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{
			"likes": bson.M{
				"$each": []bson.M{
					{
						"user_id":       userID,
						"creation_time": time.Now(),
					},
				},
				"$position": 0,
			},
		},
	}
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (e *FeedDao) UnLike(ctx *gin.Context, userID string, feedID string) error {
	collection := e.GetDB().Collection("feed")
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
