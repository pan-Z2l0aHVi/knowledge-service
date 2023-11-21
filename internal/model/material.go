package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Material struct {
	ID         primitive.ObjectID `json:"material_id" bson:"_id"`
	Type       int                `json:"type" bson:"type"`
	URL        string             `json:"url" bson:"url"`
	Name       string             `json:"name" bson:"name"`
	UploaderID string             `json:"uploader_id" bson:"uploader_id"`
}
