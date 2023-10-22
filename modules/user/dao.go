package user

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDAO struct {
	*tools.Mongo
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

func (e *UserDAO) Create(
	ctx *gin.Context,
	nickname string,
	avatar string,
	associated int,
	githubID int,
) (User, error) {
	collection := e.GetDB().Collection("user")
	user := User{
		UserID:     primitive.NewObjectID(),
		Nickname:   nickname,
		Avatar:     avatar,
		Associated: associated,
		GithubID:   githubID,
	}
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
