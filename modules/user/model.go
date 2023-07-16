package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	UserID     primitive.ObjectID `json:"user_id" bson:"_id"`
	Associated int                `json:"associated"`
	GithubID   int                `json:"github_id"`
	Nickname   string             `json:"nickname"`
	Avatar     string             `json:"avatar"`
}
