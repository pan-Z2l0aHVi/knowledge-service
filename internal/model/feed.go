package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feed struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	CreatorID     string             `json:"creator_id" bson:"creator_id"`
	SubjectID     string             `json:"subject_id" bson:"subject_id"`
	SubjectType   string             `json:"subject_type" bson:"subject_type"`
	Likes         []Like             `json:"likes" bson:"likes"`
	LikesCount    int                `json:"likes_count" bson:"likes_count"`
	Comments      []Comment          `json:"comments" bson:"comments"`
	CommentsCount int                `json:"comments_count" bson:"comments_count"`
	CreationTime  time.Time          `json:"creation_time" bson:"creation_time"`
	UpdateTime    time.Time          `json:"update_time" bson:"update_time"`
}

type Like struct {
	UserID       string    `json:"user_id" bson:"user_id"`
	CreationTime time.Time `json:"creation_time" bson:"creation_time"`
}

type Comment struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       string             `json:"user_id" bson:"user_id"`
	Content      string             `json:"content" bson:"content"`
	CreationTime time.Time          `json:"creation_time" bson:"creation_time"`
	UpdateTime   time.Time          `json:"update_time" bson:"update_time"`
	FeedID       string             `json:"feed_id" bson:"feed_id"`
	SubComments  []SubComment       `json:"sub_comments" bson:"sub_comments"`
}

type SubComment struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	UserID         string             `json:"user_id" bson:"user_id"`
	Content        string             `json:"content" bson:"content"`
	CreationTime   time.Time          `json:"creation_time" bson:"creation_time"`
	UpdateTime     time.Time          `json:"update_time" bson:"update_time"`
	FeedID         string             `json:"feed_id" bson:"feed_id"`
	ReplyCommentID string             `json:"reply_comment_id" bson:"reply_comment_id"`
}
