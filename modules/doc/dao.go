package doc

import (
	"knowledge-base-service/tools"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocDAO struct {
	*tools.Mongo
}

func (e *DocDAO) Find(ctx *gin.Context, docID string) (Doc, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return Doc{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	res := collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return Doc{}, err
	}
	var docInfo Doc
	if err := res.Decode(&docInfo); err != nil {
		return Doc{}, err
	}
	return docInfo, nil
}

func (e *DocDAO) Create(
	ctx *gin.Context,
	authorID string,
	title string,
	content string,
	cover string,
) (Doc, error) {
	collection := e.GetDB().Collection("doc")
	now := time.Now()
	doc := Doc{
		ID:           primitive.NewObjectID(),
		AuthorID:     authorID,
		Title:        title,
		Content:      content,
		Cover:        cover,
		Public:       false,
		CreationTime: now,
		UpdateTime:   now,
	}
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return Doc{}, err
	}
	return doc, nil
}

func (e *DocDAO) Update(
	ctx *gin.Context,
	docID string,
	title string,
	content string,
	cover string,
	public bool,
) (Doc, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return Doc{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "title", Value: title},
		{Key: "content", Value: content},
		{Key: "cover", Value: cover},
		{Key: "public", Value: public},
		{Key: "update_time", Value: time.Now()},
	}}}
	var doc Doc
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		return Doc{}, err
	}
	return doc, nil
}

func (e *DocDAO) Delete(ctx *gin.Context, docIDs []string) error {
	collection := e.GetDB().Collection("doc")
	var objIDs []primitive.ObjectID
	for _, docID := range docIDs {
		id, err := primitive.ObjectIDFromHex(docID)
		if err != nil {
			return err
		}
		objIDs = append(objIDs, id)
	}
	filter := bson.D{{Key: "_id", Value: bson.M{"$in": objIDs}}}
	if _, err := collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}

func (e *DocDAO) FindDocsByAuthor(ctx *gin.Context, authorID string, page int, pageSize int) ([]Doc, error) {
	collection := e.GetDB().Collection("doc")
	filter := bson.D{{Key: "author_id", Value: authorID}}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}
	var docs []Doc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	if docs == nil {
		docs = []Doc{}
	}
	return docs, nil
}

func (e *DocDAO) FindCountByAuthor(ctx *gin.Context, authorID string) (int64, error) {
	collection := e.GetDB().Collection("doc")
	filter := bson.D{{Key: "author_id", Value: authorID}}
	return collection.CountDocuments(ctx, filter)
}
