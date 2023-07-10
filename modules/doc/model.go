package doc

import "go.mongodb.org/mongo-driver/bson/primitive"

type Doc struct {
	ID       primitive.ObjectID `json:"doc_id" bson:"_id"`
	Content  string             `json:"content"`
	AuthorID string             `json:"author_id"`
}
