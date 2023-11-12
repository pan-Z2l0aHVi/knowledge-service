package space

import (
	"knowledge-base-service/tools"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SpaceDAO struct {
	*tools.Mongo
}

func (e *SpaceDAO) Find(ctx *gin.Context, spaceID string) (Space, error) {
	collection := e.GetDB().Collection("space")
	objID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return Space{}, err
	}
	filter := bson.M{"_id": objID}
	var spaceInfo Space
	if err := collection.FindOne(ctx, filter).Decode(&spaceInfo); err != nil {
		return Space{}, err
	}
	return spaceInfo, nil
}

func (e *SpaceDAO) Create(
	ctx *gin.Context,
	ownerID string,
	name string,
	desc string,
) (Space, error) {
	collection := e.GetDB().Collection("space")
	now := time.Now()
	space := Space{
		ID:           primitive.NewObjectID(),
		OwnerID:      ownerID,
		Name:         name,
		Desc:         desc,
		CreationTime: now,
		UpdateTime:   now,
	}
	_, err := collection.InsertOne(ctx, space)
	if err != nil {
		return Space{}, err
	}
	return space, nil
}

func (e *SpaceDAO) Update(
	ctx *gin.Context,
	spaceID string,
	name *string,
	desc *string,
) (Space, error) {
	collection := e.GetDB().Collection("space")
	objID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return Space{}, err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"update_time": time.Now()}}

	if name != nil {
		update["$set"].(bson.M)["name"] = name
	}
	if desc != nil {
		update["$set"].(bson.M)["desc"] = desc
	}
	var space Space
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&space); err != nil {
		return Space{}, err
	}
	return space, nil
}

func (e *SpaceDAO) Delete(ctx *gin.Context, spaceIDS []string) error {
	collection := e.GetDB().Collection("space")
	var objIDs []primitive.ObjectID
	for _, spaceID := range spaceIDS {
		id, err := primitive.ObjectIDFromHex(spaceID)
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

func (e *SpaceDAO) FindSpaces(ctx *gin.Context,
	page int,
	pageSize int,
	ownerID string,
	keywords string,
	sortBy string,
	asc int,
) ([]Space, error) {
	collection := e.GetDB().Collection("space")
	filter := bson.M{}
	if keywords != "" {
		escapedKeyword := regexp.QuoteMeta(keywords)
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": escapedKeyword, "$options": "i"}},
			{"desc": bson.M{"$regex": escapedKeyword, "$options": "i"}},
		}
	}
	if ownerID != "" {
		filter["owner_id"] = ownerID
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
	var spaces []Space
	if err := cursor.All(ctx, &spaces); err != nil {
		return nil, err
	}
	if spaces == nil {
		spaces = []Space{}
	}
	return spaces, nil
}
