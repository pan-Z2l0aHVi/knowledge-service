package dao

import (
	"knowledge-service/pkg/tools"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CommonDAO struct {
	*tools.Mongo
}

func (e *CommonDAO) InsertReport(ctx *gin.Context, jsonData []interface{}) error {
	collection := e.GetDB().Collection("fe_report")
	_, err := collection.InsertMany(ctx, jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (e *CommonDAO) FindPVCount(ctx *gin.Context, startTimestamp int64, endTimestamp int64) (int64, error) {
	collection := e.GetDB().Collection("fe_report")
	if endTimestamp == 0 {
		endTimestamp = time.Now().UnixMilli()
	}
	filter := bson.M{
		"date": bson.M{
			"$gte": time.UnixMilli(startTimestamp),
			"$lte": time.UnixMilli(endTimestamp),
		},
	}
	return collection.CountDocuments(ctx, filter)
}

func (e *CommonDAO) FindUVCount(ctx *gin.Context, startTimestamp, endTimestamp int64) (int64, error) {
	collection := e.GetDB().Collection("fe_report")
	if endTimestamp == 0 {
		endTimestamp = time.Now().UnixMilli()
	}
	filter := bson.M{
		"date": bson.M{
			"$gte": time.UnixMilli(startTimestamp),
			"$lte": time.UnixMilli(endTimestamp),
		},
	}
	pipeline := bson.A{
		bson.M{"$match": filter},
		bson.M{
			"$group": bson.M{
				"_id": "$ip",
			},
		},
		bson.M{"$count": "count"},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	var result struct {
		Count int64 `bson:"count"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}
	return result.Count, nil
}
