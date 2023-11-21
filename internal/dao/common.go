package dao

import (
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
)

type CommonDao struct {
	*tools.Mongo
}

func (e *CommonDao) InsertReport(ctx *gin.Context, jsonData []interface{}) error {
	collection := e.GetDB().Collection("fe_report")
	_, err := collection.InsertMany(ctx, jsonData)
	if err != nil {
		return err
	}
	return nil
}
