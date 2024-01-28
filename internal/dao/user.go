package dao

import (
	"encoding/json"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/model"
	"knowledge-service/pkg/tools"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDAO struct {
	*tools.Mongo
	*tools.Redis
}

func (e *UserDAO) FindByUserID(ctx *gin.Context, userID string) (model.User, error) {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return model.User{}, err
	}
	filter := bson.M{"_id": objID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return model.User{}, err
	}
	var user model.User
	if err := res.Decode(&user); err != nil {
		return model.User{}, err
	}
	if user.FollowedUserIDs == nil {
		user.FollowedUserIDs = []string{}
	}
	if user.CollectedFeedIDs == nil {
		user.CollectedFeedIDs = []string{}
	}
	if user.CollectedWallpapers == nil {
		user.CollectedWallpapers = []model.Wallpaper{}
	}
	return user, nil
}

func (e *UserDAO) FindByGithubID(ctx *gin.Context, githubID int) (model.User, error) {
	collection := e.GetDB().Collection("user")
	filter := bson.M{"github_id": githubID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return model.User{}, err
	}
	var user model.User
	if err := res.Decode(&user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (e *UserDAO) FindByWeChatID(ctx *gin.Context, wechatID string) (model.User, error) {
	collection := e.GetDB().Collection("user")
	filter := bson.M{"wechat_id": wechatID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return model.User{}, err
	}
	var user model.User
	if err := res.Decode(&user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (e *UserDAO) Create(
	ctx *gin.Context,
	nickname string,
	avatar string,
	associated int,
	githubID int,
	wechatID string,
) (model.User, error) {
	collection := e.GetDB().Collection("user")
	user := model.User{
		UserID:       primitive.NewObjectID(),
		Nickname:     nickname,
		Avatar:       avatar,
		Associated:   associated,
		GithubID:     githubID,
		WeChatID:     wechatID,
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
	}
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (e *UserDAO) Update(ctx *gin.Context, userID string, nickname *string, avatar *string) (model.User, error) {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return model.User{}, err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"update_time": time.Now()}}

	if nickname != nil {
		update["$set"].(bson.M)["nickname"] = nickname
	}
	if avatar != nil {
		update["$set"].(bson.M)["avatar"] = avatar
	}
	var user model.User
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (e *UserDAO) SetTempUserID(tempUserID string, userInfo entity.WeChatUserInfo) error {
	userJSON, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}
	rds := e.GetRDS()
	return rds.Set(tempUserID, userJSON, 300*time.Second).Err()
}

func (e *UserDAO) GetTempUserIDUserInfo(tempUserID string) (entity.WeChatUserInfo, error) {
	rds := e.GetRDS()
	resJSON, err := rds.Get(tempUserID).Result()
	if err == redis.Nil {
		return entity.WeChatUserInfo{}, err
	} else if err != nil {
		return entity.WeChatUserInfo{}, err
	}
	var userInfo entity.WeChatUserInfo
	err = json.Unmarshal([]byte(resJSON), &userInfo)
	if err != nil {
		return entity.WeChatUserInfo{}, err
	}
	return userInfo, err
}

func (e *UserDAO) AddFeedIDToCollection(ctx *gin.Context, userID string, feedID string) error {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$push": bson.M{"collected_feed_ids": feedID}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (e *UserDAO) RemoveFeedIDFromCollection(ctx *gin.Context, userID string, feedID string) error {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$pull": bson.M{"collected_feed_ids": feedID}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (e *UserDAO) FindCollectedFeedIDs(ctx *gin.Context, userID string) ([]string, error) {
	user, err := e.FindByUserID(ctx, userID)
	if err != nil {
		return []string{}, err
	}
	return user.CollectedFeedIDs, nil
}

func (e *UserDAO) AddFollowedUserID(ctx *gin.Context, userID string, targetUserID string) error {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$push": bson.M{"followed_user_ids": targetUserID}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (e *UserDAO) RemoveFollowedUserID(ctx *gin.Context, userID string, targetUserID string) error {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$pull": bson.M{"followed_user_ids": targetUserID}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (e *UserDAO) FindFollowedUserIDs(ctx *gin.Context, userID string) ([]string, error) {
	user, err := e.FindByUserID(ctx, userID)
	if err != nil {
		return []string{}, err
	}
	return user.FollowedUserIDs, nil
}

func (e *UserDAO) AddWallpaperToCollection(ctx *gin.Context, userID string, wallpaper model.Wallpaper) error {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$push": bson.M{"collected_wallpapers": wallpaper}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (e *UserDAO) RemoveWallpaperFromCollection(ctx *gin.Context, userID string, wallpaperID string) error {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$pull": bson.M{"collected_wallpapers": bson.M{"id": wallpaperID}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (e *UserDAO) FindCollectedWallpapers(ctx *gin.Context, userID string) ([]entity.WallpaperItem, error) {
	user, err := e.FindByUserID(ctx, userID)
	if err != nil {
		return []entity.WallpaperItem{}, err
	}
	wallpapers := []entity.WallpaperItem{}
	for _, wallpaper := range user.CollectedWallpapers {
		wallpapers = append(wallpapers, entity.WallpaperItem{
			Wallpaper: wallpaper,
			Collected: true,
		})
	}
	return wallpapers, nil
}
