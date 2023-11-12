package doc

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Doc struct {
	ID           primitive.ObjectID `json:"doc_id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Content      string             `json:"content" bson:"content"`
	Summary      string             `json:"summary" bson:"summary"`
	AuthorID     string             `json:"author_id" bson:"author_id"`
	SpaceID      string             `json:"space_id" bson:"space_id"`
	Cover        string             `json:"cover" bson:"cover"`
	Public       bool               `json:"public" bson:"public"`
	CreationTime time.Time          `json:"creation_time" bson:"creation_time"`
	UpdateTime   time.Time          `json:"update_time" bson:"update_time"`
}

type Author struct {
	UserID   primitive.ObjectID `json:"user_id" bson:"_id"`
	Nickname string             `json:"nickname" bson:"nickname"`
	Avatar   string             `json:"avatar" bson:"avatar"`
}

type Draft struct {
	Content      string    `json:"content" bson:"content"`
	CreationTime time.Time `json:"creation_time" bson:"creation_time"`
}
