package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"
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

func (e *User) GetProfile(ctx *gin.Context) {
	var query GetProfileQuery
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
	dao := UserDAO{}
	user, err := dao.FindByUserID(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, user)
}

func (e *User) UpdateProfile(ctx *gin.Context) {

}

func (e *User) Login(ctx *gin.Context) {
	var payload LoginPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	switch payload.Type {
	case TypeGithub:
		res, err := githubLogin(ctx, payload.Code)
		if err != nil {
			tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
			return
		}
		tools.RespSuccess(ctx, res)
	default:
		tools.RespFail(ctx, consts.FailCode, "未知登录类型", nil)
	}
}

func githubLogin(ctx *gin.Context, code string) (LoginRes, error) {
	tokenResp, err := getGitHubToken(code)
	if err != nil {
		return LoginRes{}, err
	}
	githubProfile, err := getGithubProfile(tokenResp.AccessToken)
	if err != nil {
		return LoginRes{}, err
	}
	githubID := githubProfile.ID
	dao := UserDAO{}
	user, err := dao.FindByGithubID(ctx, githubID)
	if err != nil {
		createdUser, err := dao.Create(
			ctx,
			githubProfile.Name,
			githubProfile.AvatarURL,
			TypeGithub,
			githubID,
			"",
		)
		if err != nil {
			return LoginRes{}, err
		}
		user = createdUser
	}
	token, err := tools.CreateToken(user.UserID.Hex())
	if err != nil {
		return LoginRes{}, err
	}
	res := LoginRes{
		User:  user,
		Token: token,
	}
	return res, nil
}

func getGitHubToken(code string) (GitHubTokenSuccessResp, error) {
	params := GitHubTokenPayload{
		ClientID:     GithubClientID,
		ClientSecret: GithubClientSecret,
		Code:         code,
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	req, err := http.NewRequest("POST", GithubAccessTokenURL, bytes.NewBuffer(jsonParams))
	if err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	tokenResp := GitHubTokenResp{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	successResp := GitHubTokenSuccessResp{
		AccessToken: tokenResp.AccessToken,
		Scope:       tokenResp.Scope,
		TokenType:   tokenResp.TokenType,
	}
	if len(tokenResp.Error) > 0 {
		return GitHubTokenSuccessResp{}, errors.New(tokenResp.ErrorDescription)
	}
	return successResp, nil
}

func getGithubProfile(token string) (GithubProfileResp, error) {
	req, err := http.NewRequest("GET", GithubUserAPI, nil)
	if err != nil {
		return GithubProfileResp{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return GithubProfileResp{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GithubProfileResp{}, err
	}
	profile := GithubProfileResp{}
	if err := json.Unmarshal(body, &profile); err != nil {
		return GithubProfileResp{}, err
	}
	return profile, nil
}

func (e *User) GetYDQRCode(ctx *gin.Context) {
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
	var result YDWechatQRCodeResp
	if err = json.Unmarshal(body, &result); err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	if result.Code != 0 {
		tools.RespFail(ctx, consts.FailCode, result.Msg, nil)
		return
	}
	dao := UserDAO{}
	dao.SetTempUserID(result.Data.TempUserID, 0)
	res := GetYDQRCodeResp{
		TempUserID: result.Data.TempUserID,
		QRCodeURL:  result.Data.QRURL,
	}
	tools.RespSuccess(ctx, res)
}

func (e *User) GetYDLoginStatus(ctx *gin.Context) {
	var query GetYDLoginStatusQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := UserDAO{}
	hasLogin, err := dao.GetTempUserIDLoginStatus(query.TempUserID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	if hasLogin == 1 {
		tools.RespSuccess(ctx, GetYDLoginStatusResp{
			HasLogin: true,
		})
	} else {
		tools.RespSuccess(ctx, GetYDLoginStatusResp{
			HasLogin: false,
		})
	}
}

func (e *User) YDCallback(ctx *gin.Context) {
	var payload YDCallbackPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := UserDAO{}
	_, err := dao.GetTempUserIDLoginStatus(payload.TempUserId)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
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
	res, err := wechatLogin(
		ctx,
		payload.WxMaUserInfo.OpenID,
		payload.WxMaUserInfo.Nickname,
		payload.WxMaUserInfo.AvatarUrl,
	)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	dao.SetTempUserID(payload.TempUserId, 1)
	tools.RespSuccess(ctx, res)
}

func wechatLogin(
	ctx *gin.Context,
	wechatID string,
	nickname string,
	avatar string,
) (LoginRes, error) {
	dao := UserDAO{}
	user, err := dao.FindByWeChatID(ctx, wechatID)
	if err != nil {
		createdUser, err := dao.Create(
			ctx,
			nickname,
			avatar,
			TypeWechat,
			0,
			wechatID,
		)
		if err != nil {
			return LoginRes{}, err
		}
		user = createdUser
	}
	token, err := tools.CreateToken(user.UserID.Hex())
	if err != nil {
		return LoginRes{}, err
	}
	res := LoginRes{
		User:  user,
		Token: token,
	}
	return res, nil
}
