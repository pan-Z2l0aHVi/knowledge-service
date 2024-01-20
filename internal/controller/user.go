package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/model"
	"knowledge-service/internal/service"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct{}

func (e *UserController) GetProfile(ctx *gin.Context) {
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	user, err := userD.FindByUserID(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.Profile{
		UserID:       user.UserID.Hex(),
		Associated:   user.Associated,
		GithubID:     user.GithubID,
		WeChatID:     user.WeChatID,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		CreationTime: user.CreationTime,
		UpdateTime:   user.UpdateTime,
	}
	tools.RespSuccess(ctx, res)
}

func (e *UserController) UpdateProfile(ctx *gin.Context) {
	var payload entity.UpdateProfilePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	user, err := userD.Update(ctx, userID, payload.Nickname, payload.Avatar)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.Profile{
		UserID:       user.UserID.Hex(),
		Associated:   user.Associated,
		GithubID:     user.GithubID,
		WeChatID:     user.WeChatID,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		CreationTime: user.CreationTime,
		UpdateTime:   user.UpdateTime,
	}
	tools.RespSuccess(ctx, res)
}

func (e *UserController) GetUserInfo(ctx *gin.Context) {
	var query entity.GetUserInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	userD := dao.UserDAO{}
	user, err := userD.FindByUserID(ctx, query.UserID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	collected := false
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	if userID != "" {
		self, err := userD.FindByUserID(ctx, userID)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		for _, followedUser := range self.FollowedUserIDs {
			if followedUser == query.UserID {
				collected = true
				break
			}
		}
	}
	res := entity.UserItem{
		Profile: entity.Profile{
			UserID:       user.UserID.Hex(),
			Associated:   user.Associated,
			GithubID:     user.GithubID,
			WeChatID:     user.WeChatID,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			CreationTime: user.CreationTime,
			UpdateTime:   user.UpdateTime,
		},
		Collected: collected,
	}
	tools.RespSuccess(ctx, res)
}

func (e *UserController) Login(ctx *gin.Context) {
	var payload entity.LoginPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var user model.User
	switch payload.Type {
	case consts.TypeGithub:
		userS := service.UserService{}
		userInfo, err := userS.GithubLogin(ctx, payload.Code)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		user = userInfo

	case consts.TypeWechat:
		userS := service.UserService{}
		userInfo, err := userS.WechatLogin(ctx, payload.Code)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		user = userInfo

	default:
		tools.RespFail(ctx, consts.Fail, "未知登录类型", nil)
		return
	}
	token, err := tools.CreateToken(user.UserID.Hex())
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.LoginRes{
		User:  user,
		Token: token,
	}
	tools.RespSuccess(ctx, res)
}

func (e *UserController) GetYDQRCode(ctx *gin.Context) {
	v := url.Values{}
	v.Set("secret", consts.YDSecret)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(consts.YDAPI + "/wxLogin/tempUserId?" + v.Encode())
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	var result entity.YDWechatQRCodeResp
	if err = json.Unmarshal(body, &result); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	if result.Code != 0 {
		tools.RespFail(ctx, consts.Fail, result.Msg, nil)
		return
	}
	userD := dao.UserDAO{}
	userD.SetTempUserID(result.Data.TempUserID, entity.WeChatUserInfo{})
	res := entity.GetYDQRCodeResp{
		TempUserID: result.Data.TempUserID,
		QRCodeURL:  result.Data.QRURL,
	}
	tools.RespSuccess(ctx, res)
}

func (e *UserController) GetYDLoginStatus(ctx *gin.Context) {
	var query entity.GetYDLoginStatusQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	userD := dao.UserDAO{}
	userInfo, err := userD.GetTempUserIDUserInfo(query.TempUserID)
	if err != nil {
		tools.RespSuccess(ctx, entity.GetYDLoginStatusResp{
			HasLogin: false,
		})
		return
	}
	tools.RespSuccess(ctx, entity.GetYDLoginStatusResp{
		HasLogin: userInfo.OpenID != "",
	})
}

func (e *UserController) YDCallback(ctx *gin.Context) {
	var payload entity.YDCallbackPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	if payload.ScanSuccess {
		fmt.Println("用户扫码成功")
		return
	}
	if payload.CancelLogin {
		fmt.Println("用户取消了扫码")
		return
	}
	fmt.Println("登录成功回调payload", payload)
	userD := dao.UserDAO{}
	userD.SetTempUserID(payload.TempUserId, payload.WxMaUserInfo)
	tools.RespSuccess(ctx, nil)
}

func (e *UserController) GetCollectedFeeds(ctx *gin.Context) {
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	feedIDs, err := userD.FindCollectedFeedIDs(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	feedD := dao.FeedDao{}
	feedS := service.FeedService{}
	feeds := []model.Feed{}
	for _, feedID := range feedIDs {
		feed, err := feedD.Find(ctx, feedID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			}
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		feeds = append(feeds, feed)
	}
	feedList, err := feedS.FormatFeedList(ctx, feeds, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, feedList)
}

func (e *UserController) CollectFeed(ctx *gin.Context) {
	var payload entity.AddFeedToCollectionPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	feedIDs, err := userD.FindCollectedFeedIDs(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	for _, feedID := range feedIDs {
		if feedID == payload.FeedID {
			tools.RespSuccess(ctx, nil)
			return
		}
	}
	if err := userD.AddFeedIDToCollection(ctx, userID, payload.FeedID); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}

func (e *UserController) CancelCollectFeed(ctx *gin.Context) {
	var payload entity.RemoveFeedFromCollectionPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	if err := userD.RemoveFeedIDFromCollection(ctx, userID, payload.FeedID); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}

func (e *UserController) GetFollowedUsers(ctx *gin.Context) {
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	followedUserIDs, err := userD.FindFollowedUserIDs(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	users := []entity.UserItem{}
	for _, followedUserID := range followedUserIDs {
		user, err := userD.FindByUserID(ctx, followedUserID)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		users = append(users, entity.UserItem{
			Profile: entity.Profile{
				UserID:       user.UserID.Hex(),
				Associated:   user.Associated,
				GithubID:     user.GithubID,
				WeChatID:     user.WeChatID,
				Nickname:     user.Nickname,
				Avatar:       user.Avatar,
				UpdateTime:   user.UpdateTime,
				CreationTime: user.CreationTime,
			},
			Collected: true,
		})
	}
	tools.RespSuccess(ctx, users)
}

func (e *UserController) FollowUser(ctx *gin.Context) {
	var payload entity.FollowUserPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	userIDs, err := userD.FindFollowedUserIDs(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	for _, followedUserID := range userIDs {
		if followedUserID == payload.UserID {
			tools.RespSuccess(ctx, nil)
			return
		}
	}
	if err := userD.AddFollowedUserID(ctx, userID, payload.UserID); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	user, err := userD.FindByUserID(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, user)
}

func (e *UserController) UnfollowUser(ctx *gin.Context) {
	var payload entity.UnfollowUserPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	if err := userD.RemoveFollowedUserID(ctx, userID, payload.UserID); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	user, err := userD.FindByUserID(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, user)
}

func (e *UserController) GetCollectedWallpapers(ctx *gin.Context) {
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	wallpapers, err := userD.FindCollectedWallpapers(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, wallpapers)
}

func (e *UserController) CollectWallpaper(ctx *gin.Context) {
	var payload entity.AddWallpaperToCollectionPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	user, err := userD.FindByUserID(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	for _, wallpaper := range user.CollectedWallpapers {
		if wallpaper.WallpaperID == payload.Wallpaper.WallpaperID {
			tools.RespSuccess(ctx, nil)
			return
		}
	}
	if err := userD.AddWallpaperToCollection(ctx, userID, payload.Wallpaper); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}

func (e *UserController) CancelCollectWallpaper(ctx *gin.Context) {
	var payload entity.RemoveWallpaperFromCollectionPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	userD := dao.UserDAO{}
	if err := userD.RemoveWallpaperFromCollection(ctx, userID, payload.WallpaperID); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}
