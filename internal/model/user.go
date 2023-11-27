package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID              primitive.ObjectID `json:"user_id" bson:"_id"`
	Associated          int                `json:"associated" bson:"associated"`
	GithubID            int                `json:"github_id" bson:"github_id"`
	WeChatID            string             `json:"wechat_id" bson:"wechat_id"`
	Nickname            string             `json:"nickname" bson:"nickname"`
	Avatar              string             `json:"avatar" bson:"avatar"`
	CreationTime        time.Time          `json:"creation_time" bson:"creation_time"`
	UpdateTime          time.Time          `json:"update_time" bson:"update_time"`
	CollectedFeedIDs    []string           `json:"collected_feed_ids" bson:"collected_feed_ids"`
	FollowedUserIDs     []string           `json:"followed_user_ids" bson:"followed_user_ids"`
	CollectedWallpapers []Wallpaper        `json:"collected_wallpapers" bson:"collected_wallpapers"`
}
