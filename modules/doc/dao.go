package doc

import (
	"knowledge-base-service/tools"
	"net/url"
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
	filter := bson.M{"_id": objID}
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
	title *string,
	content *string,
	summary *string,
	cover *string,
	public *bool,
) (Doc, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return Doc{}, err
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
	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	if _, err := collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}

func (e *DocDAO) FindDocs(ctx *gin.Context,
	page int,
	pageSize int,
	authorID string,
	keywords string,
	sortBy string,
	asc int,
) ([]Doc, error) {
	collection := e.GetDB().Collection("doc")
	filter := bson.M{}
	if keywords != "" {
		escapedKeyword := url.QueryEscape(keywords)
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": escapedKeyword, "$options": "i"}},
			{"summary": bson.M{"$regex": escapedKeyword, "$options": "i"}},
		}
	}
	if authorID != "" {
		filter["author_id"] = authorID
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

func (e *DocDAO) FindDraftsByDoc(ctx *gin.Context, docID string, page int, pageSize int) ([]Draft, error) {
	collection := e.GetDB().Collection("doc")
	objID, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return []Draft{}, err
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
	filter := bson.M{"_id": objID}
	newDraft := Draft{
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
		return Draft{}, err
	}
	return newDraft, nil
}
