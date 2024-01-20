package service

import (
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/model"

	"github.com/gin-gonic/gin"
)

type DocService struct{}

func (e *DocService) FormatDoc(
	ctx *gin.Context,
	doc model.Doc,
) (entity.DocInfo, error) {
	userD := dao.UserDAO{}
	author, err := userD.FindByUserID(ctx, doc.AuthorID)
	if err != nil {
		return entity.DocInfo{}, err
	}
	return entity.DocInfo{
		Doc: doc,
		Author: entity.Author{
			ID:       author.UserID.Hex(),
			Nickname: author.Nickname,
			Avatar:   author.Avatar,
		},
	}, nil
}
