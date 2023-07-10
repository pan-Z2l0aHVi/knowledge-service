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

func (e *UserDAO) find(ctx *gin.Context, userID string) (User, error) {
	collection := e.GetDB().Collection("user")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return User{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return User{}, err
	}

	var docInfo User
	if err := res.Decode(&docInfo); err != nil {
		return User{}, err
	}
	return docInfo, nil
}
