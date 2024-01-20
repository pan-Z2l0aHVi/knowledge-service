package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Space struct {
	ID           primitive.ObjectID `json:"space_id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Desc         string             `json:"desc" bson:"desc"`
	OwnerID      string             `json:"owner_id" bson:"owner_id"`
	CreationTime time.Time          `json:"creation_time" bson:"creation_time"`
	UpdateTime   time.Time          `json:"update_time" bson:"update_time"`
}
