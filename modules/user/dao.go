package user

import (
	"encoding/json"
	"knowledge-base-service/tools"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDAO struct {
	*tools.Mongo
	*tools.Redis
}

func (e *UserDAO) FindByUserID(ctx *gin.Context, userID string) (User, error) {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return User{}, err
	}
	filter := bson.M{"_id": objID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return User{}, err
	}
	var user User
	if err := res.Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

func (e *UserDAO) FindByGithubID(ctx *gin.Context, githubID int) (User, error) {
	collection := e.GetDB().Collection("user")
	filter := bson.M{"github_id": githubID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return User{}, err
	}
	var user User
	if err := res.Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

func (e *UserDAO) FindByWeChatID(ctx *gin.Context, wechatID string) (User, error) {
	collection := e.GetDB().Collection("user")
	filter := bson.M{"wechat_id": wechatID}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return User{}, err
	}
	var user User
	if err := res.Decode(&user); err != nil {
		return User{}, err
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
) (User, error) {
	collection := e.GetDB().Collection("user")
	user := User{
		UserID:       primitive.NewObjectID(),
		Nickname:     nickname,
		Avatar:       avatar,
		Associated:   associated,
		GithubID:     githubID,
		WeChatID:     wechatID,
		CreationTime: time.Now(),
	}
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (e *UserDAO) SetTempUserID(tempUserID string, userInfo WeChatUserInfo) error {
	userJSON, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}
	rds := e.GetRDS()
	return rds.Set(tempUserID, userJSON, 300*time.Second).Err()
}

func (e *UserDAO) GetTempUserIDUserInfo(tempUserID string) (WeChatUserInfo, error) {
	rds := e.GetRDS()
	resJSON, err := rds.Get(tempUserID).Result()
	if err == redis.Nil {
		return WeChatUserInfo{}, err
	} else if err != nil {
		return WeChatUserInfo{}, err
	}
	var userInfo WeChatUserInfo
	err = json.Unmarshal([]byte(resJSON), &userInfo)
	if err != nil {
		return WeChatUserInfo{}, err
	}
	return userInfo, err
}
