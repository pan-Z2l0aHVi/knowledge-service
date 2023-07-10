package material

import "go.mongodb.org/mongo-driver/bson/primitive"

type Material struct {
	ID         primitive.ObjectID `json:"material_id" bson:"_id"`
	Type       int                `json:"type"`
	URL        string             `json:"content"`
	Name       string             `json:"name"`
	UploaderID string             `json:"uploader_id"`
}
