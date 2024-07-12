package controller

import (
	"encoding/json"
	"errors"
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

// @Summary 获取本人用户信息
// @Description 校验登录态
// @Produce json
// @Success 200 {object} entity.Profile "ok" "本人用户信息"
// @Router /user/profile [get]
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

// @Summary 更新本人用户信息
// @Description 校验登录态
// @Produce json
// @Param request body entity.UpdateProfilePayload true "昵称和头像"
// @Success 200 {object} entity.Profile "ok" "更新后的本人用户信息"
// @Router /user/profile [post]
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

// @Summary 获取其他人用户信息
// @Description 使用但不校验登录态
// @Produce json
// @Param user_id query string true "用户ID"
// @Success 200 {object} entity.UserItem "ok" "用户信息"
// @Router /user/user_info [get]
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

// @Summary 用户登录
// @Description 微信授权登录｜Github授权登录
// @Produce json
// @Param request body entity.LoginPayload true "登录参数"
// @Success 200 {object} entity.LoginRes "ok" "Token和用户信息"
// @Router /user/sign_in [post]
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

// @Summary 获取易登二维码
// @Description 用于微信扫码登录
// @Produce json
// @Success 200 {object} entity.GetYDQRCodeResp "ok" "二维码链接和易登临时ID"
// @Router /user/yd_qrcode [get]
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

// @Summary 检查易登状态
// @Description 用于微信扫码登录
// @Produce json
// @Param temp_user_id query string true "易登临时ID"
// @Success 200 {object} entity.GetYDLoginStatusResp "ok" "是否成功登录"
// @Router /user/yd_login_status [get]
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

// @Summary 易登回调
// @Description 提供给易登侧调用，获取用户扫码状态
// @Produce json
// @Param request body entity.YDCallbackPayload true "登录参数"
// @Success 200 "ok"
// @Router /user/yd_callback [post]
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

// @Summary 获取收藏的动态列表
// @Description 校验登录态
// @Produce json
// @Success 200 {array} entity.FeedInfo "ok" "收藏的动态列表"
// @Router /user/collected_feeds [get]
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
	feedD := dao.FeedDAO{}
	feedS := service.FeedService{}
	feeds := []model.Feed{}
	for _, feedID := range feedIDs {
		feed, err := feedD.Find(ctx, feedID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
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

// @Summary 收藏动态
// @Description 校验登录态
// @Produce json
// @Param request body entity.AddFeedToCollectionPayload true "动态ID"
// @Success 200 "ok"
// @Router /user/collect_feed [post]
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

// @Summary 取消收藏动态
// @Description 校验登录态
// @Produce json
// @Param request body entity.RemoveFeedFromCollectionPayload true "动态ID"
// @Success 200 "ok"
// @Router /user/cancel_collect_feed [post]
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

// @Summary 获取关注的用户列表
// @Description 校验登录态
// @Produce json
// @Success 200 {array} entity.UserItem "ok" "关注的用户列表"
// @Router /user/followed_users [get]
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

// @Summary 关注用户
// @Description 校验登录态
// @Produce json
// @Param request body entity.FollowUserPayload true "用户ID"
// @Success 200 "ok"
// @Router /user/follow_user [post]
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
	tools.RespSuccess(ctx, nil)
}

// @Summary 取关用户
// @Description 校验登录态
// @Produce json
// @Param request body entity.UnfollowUserPayload true "用户ID"
// @Success 200 "ok"
// @Router /user/unfollow_user [post]
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
	tools.RespSuccess(ctx, nil)
}

// @Summary 获取收藏的壁纸列表
// @Description 校验登录态
// @Produce json
// @Success 200 {array} entity.WallpaperItem "ok" "收藏的壁纸列表"
// @Router /user/collected_wallpapers [get]
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

// @Summary 收藏壁纸
// @Description 校验登录态
// @Produce json
// @Param request body entity.AddWallpaperToCollectionPayload true "壁纸数据"
// @Success 200 "ok"
// @Router /user/collect_wallpaper [post]
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

// @Summary 取消收藏壁纸
// @Description 校验登录态
// @Produce json
// @Param request body entity.RemoveWallpaperFromCollectionPayload true "壁纸数据"
// @Success 200 "ok"
// @Router /user/cancel_collect_wallpaper [post]
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
