package feed

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feed struct {
	FeedID       primitive.ObjectID `json:"feed_id" bson:"_id"`
	LikesCount   int                `json:"likes_count" bson:"likes_count"`
	Title        string             `json:"title" bson:"title"`
	Content      string             `json:"content" bson:"content"`
	AuthorID     string             `json:"author_id" bson:"author_id"`
	Cover        string             `json:"cover" bson:"cover"`
	Public       bool               `json:"public" bson:"public"`
	CreationTime time.Time          `json:"creation_time" bson:"creation_time"`
	UpdateTime   time.Time          `json:"update_time" bson:"update_time"`
}

type Author struct {
	UserID     primitive.ObjectID `json:"user_id" bson:"_id"`
	Associated int                `json:"associated" bson:"associated"`
	GithubID   int                `json:"github_id" bson:"github_id"`
	Nickname   string             `json:"nickname" bson:"nickname"`
	Avatar     string             `json:"avatar" bson:"avatar"`
}
