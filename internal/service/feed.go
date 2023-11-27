package service

import (
	"knowledge-service/internal/api"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/model"

	"github.com/gin-gonic/gin"
)

type FeedService struct{}

func (e *FeedService) FormatFeed(
	ctx *gin.Context,
	feed model.Feed,
	collectedFeedIDs []string,
) (api.FeedItem, error) {
	feedD := dao.FeedDao{}
	author, err := feedD.FindAuthorByID(ctx, feed.AuthorID)
	if err != nil {
		return api.FeedItem{}, nil
	}
	feedIDString := feed.FeedID.Hex()
	collected := false
	for _, collectedFeedID := range collectedFeedIDs {
		if collectedFeedID == feedIDString {
			collected = true
			break
		}
	}
	return api.FeedItem{
		Feed:       feed,
		AuthorInfo: author,
		Collected:  collected,
	}, nil
}
