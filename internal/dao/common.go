package dao

import (
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
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
