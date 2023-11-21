package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"knowledge-service/internal/api"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/model"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	TypeGithub           = 1
	TypeWechat           = 2
	GithubClientID       = "623037fcf1a6cb4ad6d8"
	GithubClientSecret   = "7ccd7c57dce15c44deee8760f275085afe567708"
	GithubAccessTokenURL = "https://github.com/login/oauth/access_token"
	GithubUserAPI        = "https://api.github.com/user"
	YDAPI                = "https://yd.jylt.cc/api"
	YDSecret             = "a69aca2e"
)

type UserController struct{}

func (e *UserController) GetProfile(ctx *gin.Context) {
	var query api.GetProfileQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if query.UserID != "" {
		userID = query.UserID
	} else if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	} else {
		tools.RespFail(ctx, consts.FailCode, "当前用户不存在", nil)
		return
	}
	userD := dao.UserDAO{}
	user, err := userD.FindByUserID(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, user)
}

func (e *UserController) UpdateProfile(ctx *gin.Context) {
	var payload api.UpdateProfilePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	} else {
		tools.RespFail(ctx, consts.FailCode, "当前用户不存在", nil)
		return
	}
	userD := dao.UserDAO{}
	user, err := userD.Update(ctx, userID, payload.Nickname, payload.Avatar)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, user)
}

func (e *UserController) Login(ctx *gin.Context) {
	var payload api.LoginPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	var user model.User
	switch payload.Type {
	case TypeGithub:
		userInfo, err := githubLogin(ctx, payload.Code)
		if err != nil {
			tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
			return
		}
		user = userInfo

	case TypeWechat:
		userInfo, err := wechatLogin(ctx, payload.Code)
		if err != nil {
			tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
			return
		}
		user = userInfo

	default:
		tools.RespFail(ctx, consts.FailCode, "未知登录类型", nil)
		return
	}
	token, err := tools.CreateToken(user.UserID.Hex())
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := api.LoginRes{
		User:  user,
		Token: token,
	}
	tools.RespSuccess(ctx, res)
}

func githubLogin(ctx *gin.Context, code string) (model.User, error) {
	tokenResp, err := getGitHubToken(code)
	if err != nil {
		return model.User{}, err
	}
	githubProfile, err := getGithubProfile(tokenResp.AccessToken)
	if err != nil {
		return model.User{}, err
	}
	githubID := githubProfile.ID
	userD := dao.UserDAO{}
	user, err := userD.FindByGithubID(ctx, githubID)
	if err != nil {
		createdUser, err := userD.Create(
			ctx,
			githubProfile.Name,
			githubProfile.AvatarURL,
			TypeGithub,
			githubID,
			"",
		)
		if err != nil {
			return model.User{}, err
		}
		user = createdUser
	}
	return user, nil
}

func getGitHubToken(code string) (api.GitHubTokenSuccessResp, error) {
	params := api.GitHubTokenPayload{
		ClientID:     GithubClientID,
		ClientSecret: GithubClientSecret,
		Code:         code,
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return api.GitHubTokenSuccessResp{}, err
	}
	req, err := http.NewRequest("POST", GithubAccessTokenURL, bytes.NewBuffer(jsonParams))
	if err != nil {
		return api.GitHubTokenSuccessResp{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return api.GitHubTokenSuccessResp{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.GitHubTokenSuccessResp{}, err
	}
	tokenResp := api.GitHubTokenResp{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return api.GitHubTokenSuccessResp{}, err
	}
	successResp := api.GitHubTokenSuccessResp{
		AccessToken: tokenResp.AccessToken,
		Scope:       tokenResp.Scope,
		TokenType:   tokenResp.TokenType,
	}
	if len(tokenResp.Error) > 0 {
		return api.GitHubTokenSuccessResp{}, errors.New(tokenResp.ErrorDescription)
	}
	return successResp, nil
}

func getGithubProfile(token string) (api.GithubProfileResp, error) {
	req, err := http.NewRequest("GET", GithubUserAPI, nil)
	if err != nil {
		return api.GithubProfileResp{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return api.GithubProfileResp{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.GithubProfileResp{}, err
	}
	profile := api.GithubProfileResp{}
	if err := json.Unmarshal(body, &profile); err != nil {
		return api.GithubProfileResp{}, err
	}
	return profile, nil
}

func (e *UserController) GetYDQRCode(ctx *gin.Context) {
	v := url.Values{}
	v.Set("secret", YDSecret)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(YDAPI + "/wxLogin/tempUserId?" + v.Encode())
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	var result api.YDWechatQRCodeResp
	if err = json.Unmarshal(body, &result); err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	if result.Code != 0 {
		tools.RespFail(ctx, consts.FailCode, result.Msg, nil)
		return
	}
	userD := dao.UserDAO{}
	userD.SetTempUserID(result.Data.TempUserID, api.WeChatUserInfo{})
	res := api.GetYDQRCodeResp{
		TempUserID: result.Data.TempUserID,
		QRCodeURL:  result.Data.QRURL,
	}
	tools.RespSuccess(ctx, res)
}

func (e *UserController) GetYDLoginStatus(ctx *gin.Context) {
	var query api.GetYDLoginStatusQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	userD := dao.UserDAO{}
	userInfo, err := userD.GetTempUserIDUserInfo(query.TempUserID)
	if err != nil {
		tools.RespSuccess(ctx, api.GetYDLoginStatusResp{
			HasLogin: false,
		})
		return
	}
	tools.RespSuccess(ctx, api.GetYDLoginStatusResp{
		HasLogin: userInfo.OpenID != "",
	})
}

func (e *UserController) YDCallback(ctx *gin.Context) {
	var payload api.YDCallbackPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
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

func wechatLogin(
	ctx *gin.Context,
	code string,
) (model.User, error) {
	userD := dao.UserDAO{}
	userInfo, err := userD.GetTempUserIDUserInfo(code)
	if err != nil {
		return model.User{}, err
	}
	if userInfo.OpenID == "" {
		return model.User{}, err
	}
	user, err := userD.FindByWeChatID(ctx, userInfo.OpenID)
	if err != nil {
		createdUser, err := userD.Create(
			ctx,
			userInfo.Nickname,
			userInfo.AvatarUrl,
			TypeWechat,
			0,
			userInfo.OpenID,
		)
		if err != nil {
			return model.User{}, err
		}
		user = createdUser
	}
	return user, nil
}
