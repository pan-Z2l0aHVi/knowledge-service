package dao

import (
	"encoding/json"
	"knowledge-service/internal/api"
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

func (e *UserDAO) SetTempUserID(tempUserID string, userInfo api.WeChatUserInfo) error {
	userJSON, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}
	rds := e.GetRDS()
	return rds.Set(tempUserID, userJSON, 300*time.Second).Err()
}

func (e *UserDAO) GetTempUserIDUserInfo(tempUserID string) (api.WeChatUserInfo, error) {
	rds := e.GetRDS()
	resJSON, err := rds.Get(tempUserID).Result()
	if err == redis.Nil {
		return api.WeChatUserInfo{}, err
	} else if err != nil {
		return api.WeChatUserInfo{}, err
	}
	var userInfo api.WeChatUserInfo
	err = json.Unmarshal([]byte(resJSON), &userInfo)
	if err != nil {
		return api.WeChatUserInfo{}, err
	}
	return userInfo, err
}
