package service

import (
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/model"
	"knowledge-service/pkg/consts"

	"github.com/gin-gonic/gin"
)

type DocService struct{}

func (e *DocService) FormatDoc(ctx *gin.Context, doc model.Doc) (entity.DocInfo, error) {
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

func (e *DocService) DeleteDocs(ctx *gin.Context, docIDs []string) error {
	docD := dao.DocDAO{}
	if err := docD.DeleteMany(ctx, docIDs); err != nil {
		return err
	}
	feedD := dao.FeedDAO{}
	if err := feedD.DeleteManyBySubject(ctx, docIDs, consts.DocFeed); err != nil {
		return err
	}
	feedS := FeedService{}
	if err := feedS.RemoveAllFeedListCache(); err != nil {
		return err
	}
	return nil
}
