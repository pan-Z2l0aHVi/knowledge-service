package service

import (
	"errors"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/model"
	"knowledge-service/pkg/tools"

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
	feedD := dao.FeedDAO{}
	feed, err := feedD.FindBySubject(ctx, subjectID, subjectType)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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

func (e *FeedService) FormatComments(ctx *gin.Context, comments []model.Comment) ([]entity.CommentInfo, error) {
	userD := dao.UserDAO{}
	commentList := []entity.CommentInfo{}
	for _, comment := range comments {
		commentator, err := userD.FindByUserID(ctx, comment.UserID)
		if err != nil {
			return []entity.CommentInfo{}, err
		}

		replyInfos := []entity.ReplyInfo{}
		for _, subComment := range comment.SubComments {
			subCommentUser, err := userD.FindByUserID(ctx, subComment.UserID)
			if err != nil {
				return []entity.CommentInfo{}, err
			}
			replyUser, err := userD.FindByUserID(ctx, subComment.ReplyUserID)
			if err != nil {
				return []entity.CommentInfo{}, err
			}
			replyInfos = append(replyInfos, entity.ReplyInfo{
				SubComment:       subComment,
				Commentator:      e.formatCommentator(subCommentUser),
				ReplyCommentator: e.formatCommentator(replyUser),
			})
		}

		commentList = append(commentList, entity.CommentInfo{
			Comment:     comment,
			Commentator: e.formatCommentator(commentator),
			SubComments: replyInfos,
		})
	}
	return commentList, nil
}

func (e *FeedService) FormatComment(ctx *gin.Context, comment model.Comment) (entity.CommentInfo, error) {
	commentList := append([]model.Comment{}, comment)
	list, err := e.FormatComments(ctx, commentList)
	if err != nil {
		return entity.CommentInfo{}, err
	}
	return list[0], nil
}

func (e *FeedService) FormatSubComment(ctx *gin.Context, feedID string, commentID string, subComment model.SubComment) (entity.ReplyInfo, error) {
	userD := dao.UserDAO{}
	user, err := userD.FindByUserID(ctx, subComment.UserID)
	if err != nil {
		return entity.ReplyInfo{}, err
	}
	replyUser, err := userD.FindByUserID(ctx, subComment.ReplyUserID)
	if err != nil {
		return entity.ReplyInfo{}, err
	}
	replyInfo := entity.ReplyInfo{
		SubComment:       subComment,
		Commentator:      e.formatCommentator(user),
		ReplyCommentator: e.formatCommentator(replyUser),
	}
	return replyInfo, nil
}

func (e *FeedService) formatCommentator(user model.User) entity.Commentator {
	return entity.Commentator{
		ID:       user.UserID.Hex(),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}
}

func (e *FeedService) RemoveAllFeedListCache() error {
	rdsInst := tools.Redis{}
	rds := rdsInst.GetRDS()
	pattern := "feed_list_*"
	cursor := uint64(0)
	for {
		// 扫描匹配的键
		var keys []string
		var err error
		keys, cursor, err = rds.Scan(cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		// 删除匹配的键
		if len(keys) > 0 {
			_, err = rds.Del(keys...).Result()
			if err != nil {
				return err
			}
		}
		// 如果 cursor 为 0，表示已经遍历完成
		if cursor == 0 {
			break
		}
	}
	return nil
}
