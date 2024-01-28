package dao

import (
	"knowledge-service/internal/model"
	"knowledge-service/pkg/tools"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocDAO struct {
	*tools.Mongo
}

func (e *DocDAO) Find(ctx *gin.Context, docID string) (model.Doc, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return model.Doc{}, err
	}
	filter := bson.M{"_id": objID}
	var doc model.Doc
	if err := collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		return model.Doc{}, err
	}
	return doc, nil
}

func (e *DocDAO) Create(
	ctx *gin.Context,
	authorID string,
	spaceID string,
	title string,
	content string,
	cover string,
) (model.Doc, error) {
	collection := e.GetDB().Collection("doc")
	now := time.Now()
	doc := model.Doc{
		ID:           primitive.NewObjectID(),
		AuthorID:     authorID,
		SpaceID:      spaceID,
		Title:        title,
		Content:      content,
		Cover:        cover,
		Public:       false,
		CreationTime: now,
		UpdateTime:   now,
	}
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return model.Doc{}, err
	}
	return doc, nil
}

func (e *DocDAO) Update(
	ctx *gin.Context,
	docID string,
	title *string,
	content *string,
	summary *string,
	cover *string,
	public *bool,
) (model.Doc, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return model.Doc{}, err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"update_time": time.Now()}}

	if title != nil {
		update["$set"].(bson.M)["title"] = title
	}
	if content != nil {
		update["$set"].(bson.M)["content"] = content
	}
	if summary != nil {
		update["$set"].(bson.M)["summary"] = summary
	}
	if cover != nil {
		update["$set"].(bson.M)["cover"] = cover
	}
	if public != nil {
		update["$set"].(bson.M)["public"] = *public
	}
	var doc model.Doc
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		return model.Doc{}, err
	}
	return doc, nil
}

func (e *DocDAO) DeleteMany(ctx *gin.Context, docIDs []string) error {
	collection := e.GetDB().Collection("doc")
	objIDs := []primitive.ObjectID{}
	for _, docID := range docIDs {
		id, err := primitive.ObjectIDFromHex(docID)
		if err != nil {
			return err
		}
		objIDs = append(objIDs, id)
	}
	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	if _, err := collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}

func (e *DocDAO) FindListWithTotal(ctx *gin.Context,
	page int,
	pageSize int,
	authorID string,
	spaceID string,
	keywords string,
	sortBy string,
	asc int,
	public *bool,
) ([]model.Doc, int64, error) {
	collection := e.GetDB().Collection("doc")
	filter := bson.M{}
	if keywords != "" {
		escapedKeyword := regexp.QuoteMeta(keywords)
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": escapedKeyword, "$options": "i"}},
			{"summary": bson.M{"$regex": escapedKeyword, "$options": "i"}},
		}
	}
	if authorID != "" {
		filter["author_id"] = authorID
	}
	if spaceID != "" {
		filter["space_id"] = spaceID
	}
	if public != nil {
		filter["public"] = public
	}
	sort := bson.M{}
	if sortBy != "" && asc != 0 {
		sort[sortBy] = asc
	}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  sort,
	})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	docs := []model.Doc{}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return docs, count, nil
}

func (e *DocDAO) FindDrafts(ctx *gin.Context, docID string, page int, pageSize int) ([]model.Draft, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return []model.Draft{}, err
	}
	filter := bson.M{"_id": objID}
	sort := bson.M{"creation_time": -1}
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
		return []model.Draft{}, err
	}
	defer cursor.Close(ctx)
	drafts := []model.Draft{}
	if err := cursor.All(ctx, &drafts); err != nil {
		return []model.Draft{}, err
	}
	return drafts, nil
}

func (e *DocDAO) CreateDraft(ctx *gin.Context, docID string, content string) (model.Draft, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return model.Draft{}, err
	}
	filter := bson.M{"_id": objID}
	newDraft := model.Draft{
		Content:      content,
		CreationTime: time.Now(),
	}
	update := bson.M{
		"$push": bson.M{
			"drafts": bson.M{
				"$each":     bson.A{newDraft},
				"$position": 0,
			},
		},
	}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return model.Draft{}, err
	}
	return newDraft, nil
}
