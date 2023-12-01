package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feed struct {
	FeedID       primitive.ObjectID `json:"feed_id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Content      string             `json:"content" bson:"content"`
	Summary      string             `json:"summary" bson:"summary"`
	AuthorID     string             `json:"author_id" bson:"author_id"`
	SpaceID      string             `json:"space_id" bson:"space_id"`
	Cover        string             `json:"cover" bson:"cover"`
	Public       bool               `json:"public" bson:"public"`
	CreationTime time.Time          `json:"creation_time" bson:"creation_time"`
	UpdateTime   time.Time          `json:"update_time" bson:"update_time"`
	Likes        []Like             `json:"likes" bson:"likes"`
	LikesCount   int                `json:"likes_count" bson:"likes_count"`
}

type Like struct {
	UserID       string    `json:"user_id" bson:"user_id"`
	CreationTime time.Time `json:"creation_time" bson:"creation_time"`
}
