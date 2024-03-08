package dao

import (
	"errors"
	"knowledge-service/internal/model"
	"knowledge-service/pkg/tools"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedDAO struct {
	*tools.Mongo
}

func (e *FeedDAO) Create(ctx *gin.Context, creatorID string, subjectID string, subjectType string) (model.Feed, error) {
	collection := e.GetDB().Collection("feed")
	now := time.Now()
	feed := model.Feed{
		ID:            primitive.NewObjectID(),
		CreatorID:     creatorID,
		SubjectID:     subjectID,
		SubjectType:   subjectType,
		Likes:         []model.Like{},
		LikesCount:    0,
		Comments:      []model.Comment{},
		CommentsCount: 0,
		CreationTime:  now,
		UpdateTime:    now,
	}
	_, err := collection.InsertOne(ctx, feed)
	if err != nil {
		return model.Feed{}, err
	}
	return feed, nil
}

func (e *FeedDAO) Update(ctx *gin.Context, feedID string) (model.Feed, error) {
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

func (e *FeedDAO) DeleteMany(ctx *gin.Context, feedIDs []string) error {
	collection := e.GetDB().Collection("feed")
	objIDs := []primitive.ObjectID{}
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

func (e *FeedDAO) DeleteManyBySubject(ctx *gin.Context, subjectIDs []string, subjectType string) error {
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

func (e *FeedDAO) FindListWithTotal(
	ctx *gin.Context,
	page int,
	pageSize int,
	sortBy string,
	asc int,
	authorID string,
) ([]model.Feed, int64, error) {
	collection := e.GetDB().Collection("feed")
	filter := bson.M{}
	if authorID != "" {
		filter["creator_id"] = authorID
	}
	sort := bson.M{}
	if sortBy != "" && asc != 0 {
		sort["subject."+sortBy] = asc
	}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	pipeline := bson.A{
		bson.M{
			"$addFields": bson.M{
				"subject_id_objectid": bson.M{"$toObjectId": "$subject_id"},
			},
		},
		bson.M{"$lookup": bson.M{
			"from":         "doc",
			"localField":   "subject_id_objectid",
			"foreignField": "_id",
			"as":           "subject",
		}},
		bson.M{"$unwind": "$subject"},
		bson.M{"$match": filter},
		bson.M{"$sort": sort},
		bson.M{"$skip": skip},
		bson.M{"$limit": limit},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []model.Feed{}, 0, err
	}
	defer cursor.Close(ctx)
	feeds := []model.Feed{}
	if err := cursor.All(ctx, &feeds); err != nil {
		return []model.Feed{}, 0, err
	}
	for i := range feeds {
		if feeds[i].Likes == nil {
			feeds[i].Likes = []model.Like{}
		}
		feeds[i].LikesCount = len(feeds[i].Likes)
		if feeds[i].Comments == nil {
			feeds[i].Comments = []model.Comment{}
		}
		feeds[i].CommentsCount = getCommentCount(feeds[i].Comments)
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return []model.Feed{}, 0, err
	}
	return feeds, count, nil
}

func (e *FeedDAO) Find(ctx *gin.Context, feedID string) (model.Feed, error) {
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
	if feed.Comments == nil {
		feed.Comments = []model.Comment{}
	}
	feed.CommentsCount = getCommentCount(feed.Comments)
	return feed, nil
}

func (e *FeedDAO) FindBySubject(ctx *gin.Context, subjectID string, subjectType string) (model.Feed, error) {
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
	if feed.Likes == nil {
		feed.Likes = []model.Like{}
	}
	feed.LikesCount = len(feed.Likes)
	if feed.Comments == nil {
		feed.Comments = []model.Comment{}
	}
	feed.CommentsCount = getCommentCount(feed.Comments)
	return feed, nil
}

func (e *FeedDAO) CheckLiked(ctx *gin.Context, userID string, feedID string) (bool, error) {
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

func (e *FeedDAO) Like(ctx *gin.Context, userID string, feedID string) error {
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

func (e *FeedDAO) UnLike(ctx *gin.Context, userID string, feedID string) error {
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

func (e *FeedDAO) CreateComment(
	ctx *gin.Context,
	feedID string,
	userID string,
	content string,
) (model.Comment, error) {
	collection := e.GetDB().Collection("feed")
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return model.Comment{}, err
	}
	newComment := model.Comment{
		ID:           primitive.NewObjectID(),
		FeedID:       feedID,
		UserID:       userID,
		Content:      content,
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		SubComments:  []model.SubComment{},
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{
			"comments": bson.M{
				"$each":     bson.A{newComment},
				"$position": 0,
			},
		},
	}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return model.Comment{}, err
	}
	return newComment, nil
}

func (e *FeedDAO) ReplyComment(
	ctx *gin.Context,
	feedID string,
	commentID string,
	replyUserID string,
	userID string,
	content string,
) (model.SubComment, error) {
	collection := e.GetDB().Collection("feed")
	feedObjID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return model.SubComment{}, err
	}
	commentObjID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return model.SubComment{}, err
	}
	subComment := model.SubComment{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		Content:      content,
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		FeedID:       feedID,
		ReplyUserID:  replyUserID,
	}
	filter := bson.M{
		"_id": feedObjID,
		"comments": bson.M{
			"$elemMatch": bson.M{
				"$or": []bson.M{
					{"_id": commentObjID},
					{"sub_comments._id": commentObjID},
				},
			},
		},
	}
	update := bson.M{
		"$push": bson.M{
			"comments.$.sub_comments": bson.M{
				"$each":     bson.A{subComment},
				"$position": 0,
			},
		},
	}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return model.SubComment{}, err
	}
	return subComment, nil
}

func (e *FeedDAO) DeleteComment(ctx *gin.Context, feedID string, commentID string, subCommentID string) error {
	collection := e.GetDB().Collection("feed")
	feedObjID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return err
	}
	commentObjID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": feedObjID,
	}
	update := bson.M{}
	if subCommentID != "" {
		subCommentObjID, err := primitive.ObjectIDFromHex(subCommentID)
		if err != nil {
			return err
		}
		filter["comments.sub_comments._id"] = subCommentObjID
		update["$pull"] = bson.M{"comments.$.sub_comments": bson.M{"_id": subCommentObjID}}
	} else {
		filter["comments._id"] = commentObjID
		update["$pull"] = bson.M{"comments": bson.M{"_id": commentObjID}}
	}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (e *FeedDAO) UpdateComment(
	ctx *gin.Context,
	feedID string,
	commentID string,
	content string,
) (model.Comment, error) {
	collection := e.GetDB().Collection("feed")
	feedObjID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return model.Comment{}, err
	}
	commentObjID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return model.Comment{}, err
	}
	filter := bson.M{
		"_id":          feedObjID,
		"comments._id": commentObjID,
	}
	update := bson.M{
		"$set": bson.M{
			"comments.$.content":     content,
			"comments.$.update_time": time.Now(),
		},
	}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return model.Comment{}, err
	}
	comment, err := e.FindComment(ctx, feedID, commentID)
	if err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (e *FeedDAO) UpdateSubComment(
	ctx *gin.Context,
	feedID string,
	commentID string,
	subCommentID string,
	newContent string,
) (model.SubComment, error) {
	collection := e.GetDB().Collection("feed")
	feedObjID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return model.SubComment{}, err
	}
	commentObjID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return model.SubComment{}, err
	}
	subCommentObjID, err := primitive.ObjectIDFromHex(subCommentID)
	if err != nil {
		return model.SubComment{}, err
	}
	filter := bson.M{
		"_id": feedObjID,
		"comments": bson.M{
			"$elemMatch": bson.M{
				"_id":              commentObjID,
				"sub_comments._id": subCommentObjID,
			},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"comments.$.sub_comments.$[subComment].content":     newContent,
			"comments.$.sub_comments.$[subComment].update_time": time.Now(),
		},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "subComment._id", Value: subCommentObjID}},
		},
	})
	if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
		return model.SubComment{}, err
	}
	subComment, err := e.FindSubComment(ctx, feedID, commentID, subCommentID)
	if err != nil {
		return model.SubComment{}, err
	}
	return subComment, nil
}

func (e *FeedDAO) FindComment(
	ctx *gin.Context,
	feedID string,
	commentID string,
) (model.Comment, error) {
	feed, err := e.Find(ctx, feedID)
	if err != nil {
		return model.Comment{}, err
	}
	for _, comment := range feed.Comments {
		if comment.ID.Hex() == commentID {
			return comment, nil
		}
	}
	return model.Comment{}, errors.New("comment not found")
}

func (e *FeedDAO) FindSubComment(
	ctx *gin.Context,
	feedID string,
	commentID string,
	subCommentID string,
) (model.SubComment, error) {
	comment, err := e.FindComment(ctx, feedID, commentID)
	if err != nil {
		return model.SubComment{}, err
	}
	if subCommentID != "" {
		for _, subComment := range comment.SubComments {
			if subComment.ID.Hex() == subCommentID {
				return subComment, nil
			}
		}
	}
	return model.SubComment{}, errors.New("sub comment not found")
}

func (e *FeedDAO) FindCommentListWithTotal(
	ctx *gin.Context,
	feedID string,
	page int,
	pageSize int,
	sortBy string,
	asc int,
) ([]model.Comment, int, error) {
	collection := e.GetDB().Collection("feed")
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return []model.Comment{}, 0, err
	}
	sortStage := bson.M{"comments." + sortBy: asc}
	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": objID}},
		bson.M{"$unwind": "$comments"},
		bson.M{"$sort": sortStage},
		bson.M{"$skip": int64((page - 1) * pageSize)},
		bson.M{"$limit": int64(pageSize)},
		bson.M{"$project": bson.M{
			"_id":           "$comments._id",
			"user_id":       "$comments.user_id",
			"content":       "$comments.content",
			"creation_time": "$comments.creation_time",
			"update_time":   "$comments.update_time",
			"feed_id":       "$comments.feed_id",
			"sub_comments":  "$comments.sub_comments",
		}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []model.Comment{}, 0, err
	}
	defer cursor.Close(ctx)
	feedComments := []model.Comment{}
	if err := cursor.All(ctx, &feedComments); err != nil {
		return []model.Comment{}, 0, err
	}
	feed, err := e.Find(ctx, feedID)
	if err != nil {
		return []model.Comment{}, 0, err
	}
	for _, comment := range feedComments {
		sort.Sort(ByUpdateTime(comment.SubComments))
	}
	return feedComments, feed.CommentsCount, nil
}

func getCommentCount(comments []model.Comment) int {
	total := len(comments)
	for _, comment := range comments {
		total += len(comment.SubComments)
	}
	return total
}

type ByUpdateTime []model.SubComment

func (a ByUpdateTime) Len() int {
	return len(a)
}
func (a ByUpdateTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByUpdateTime) Less(i, j int) bool {
	return a[i].UpdateTime.After(a[j].UpdateTime)
}
