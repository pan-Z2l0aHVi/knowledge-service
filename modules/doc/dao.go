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
	var docInfo Doc
	if err := collection.FindOne(ctx, filter).Decode(&docInfo); err != nil {
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
	public *bool,
) (Doc, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return Doc{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "update_time", Value: time.Now()},
	}}}
	if title != "" {
		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "title", Value: title})
	}
	if content != "" {
		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "content", Value: content})
	}
	if cover != "" {
		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "cover", Value: cover})
	}
	if public != nil {
		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "public", Value: *public})
	}
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
	sort := bson.D{{Key: "update_time", Value: -1}}
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  sort,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
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

func (e *DocDAO) FindDraftsByDoc(ctx *gin.Context, docID string, page int, pageSize int) ([]Draft, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return []Draft{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	sort := bson.D{{Key: "creation_time", Value: -1}}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	pipeline := []bson.M{
		{"$match": filter},
		{"$unwind": "$drafts"},
		{"$sort": sort},
		{"$skip": skip},
		{"$limit": limit},
		{
			"$project": bson.M{
				"_id":    0,
				"drafts": 1,
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$drafts",
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []Draft{}, err
	}
	defer cursor.Close(ctx)
	var drafts []Draft
	if err := cursor.All(ctx, &drafts); err != nil {
		return []Draft{}, err
	}
	if drafts == nil {
		drafts = []Draft{}
	}
	return drafts, nil
}

func (e *DocDAO) UpdateDraft(ctx *gin.Context, docID string, content string) (Draft, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return Draft{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	newDraft := Draft{
		Content:      content,
		CreationTime: time.Now(),
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "drafts", Value: bson.D{
				{Key: "$each", Value: bson.A{newDraft}},
				{Key: "$position", Value: 0},
			}},
		}},
	}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return Draft{}, err
	}
	return newDraft, nil
}
