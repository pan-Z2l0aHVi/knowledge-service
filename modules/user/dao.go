package user

import (
	"knowledge-base-service/tools"
	"time"

	"github.com/gin-gonic/gin"
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
		UserID:     primitive.NewObjectID(),
		Nickname:   nickname,
		Avatar:     avatar,
		Associated: associated,
		GithubID:   githubID,
		WeChatID:   wechatID,
	}
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (e *UserDAO) SetTempUserID(tempUserID string, hasLogin int) error {
	rds := e.GetRDS()
	return rds.Set(tempUserID, hasLogin, 300*time.Second).Err()
}

func (e *UserDAO) GetTempUserIDLoginStatus(tempUserID string) (int, error) {
	rds := e.GetRDS()
	res, err := rds.Get(tempUserID).Int()
	if res == 0 {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return res, err
}
