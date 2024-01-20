package service

import (
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedService struct{}

func (e *FeedService) formatFeedInfo(ctx *gin.Context, feed model.Feed) (entity.FeedInfo, error) {
	userD := dao.UserDAO{}
	creator, err := userD.FindByUserID(ctx, feed.CreatorID)
	if err != nil {
		return entity.FeedInfo{}, err
	}
	docD := dao.DocDAO{}
	doc, err := docD.Find(ctx, feed.SubjectID)
	if err != nil {
		return entity.FeedInfo{}, err
	}
	docInfo := entity.DocInfo{
		Doc: doc,
		Author: entity.Author{
			ID:       creator.UserID.Hex(),
			Nickname: creator.Nickname,
			Avatar:   creator.Avatar,
		},
	}
	likes := []entity.LikeInfo{}
	for _, like := range feed.Likes {
		user, err := userD.FindByUserID(ctx, like.UserID)
		if err != nil {
			return entity.FeedInfo{}, err
		}
		likes = append(likes, entity.LikeInfo{
			Like:     like,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		})
	}
	return entity.FeedInfo{
		Feed:  feed,
		Likes: likes,
		Creator: entity.Creator{
			ID:       creator.UserID.Hex(),
			Nickname: creator.Nickname,
			Avatar:   creator.Avatar,
		},
		Subject: docInfo,
	}, nil
}

func (e *FeedService) FormatFeedList(ctx *gin.Context, feeds []model.Feed, userID string) ([]entity.FeedInfo, error) {
	collectedFeedIDs := []string{}
	if userID != "" {
		userD := dao.UserDAO{}
		feedIDs, err := userD.FindCollectedFeedIDs(ctx, userID)
		if err != nil {
			return []entity.FeedInfo{}, err
		}
		collectedFeedIDs = feedIDs
	}
	feedList := []entity.FeedInfo{}
	for _, feed := range feeds {
		feedInfo, err := e.formatFeedInfo(ctx, feed)
		if err != nil {
			return []entity.FeedInfo{}, err
		}
		collected := false
		for _, collectedFeedID := range collectedFeedIDs {
			if collectedFeedID == feed.ID.Hex() {
				collected = true
				break
			}
		}
		feedInfo.Collected = collected
		feedList = append(feedList, feedInfo)
	}
	return feedList, nil
}

func (e *FeedService) FormatFeed(ctx *gin.Context, feed model.Feed, userID string) (entity.FeedInfo, error) {
	feeds := append([]model.Feed{}, feed)
	feedList, err := e.FormatFeedList(ctx, feeds, userID)
	if err != nil {
		return entity.FeedInfo{}, err
	}
	return feedList[0], nil
}

func (e *FeedService) SyncFeed(
	ctx *gin.Context,
	userID string,
	creatorID string,
	subjectID string,
	subjectType string,
) (entity.FeedInfo, error) {
	feedD := dao.FeedDao{}
	feed, err := feedD.FindBySubject(ctx, subjectID, subjectType)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newFeed, err := feedD.Create(ctx, creatorID, subjectID, subjectType)
			if err != nil {
				return entity.FeedInfo{}, err
			}
			feedInfo, err := e.FormatFeed(ctx, newFeed, userID)
			if err != nil {
				return entity.FeedInfo{}, err
			}
			return feedInfo, nil
		}
		return entity.FeedInfo{}, err
	}
	updatedFeed, err := feedD.Update(ctx, feed.ID.Hex())
	if err != nil {
		return entity.FeedInfo{}, err
	}
	feedInfo, err := e.FormatFeed(ctx, updatedFeed, userID)
	if err != nil {
		return entity.FeedInfo{}, err
	}
	return feedInfo, nil
}
