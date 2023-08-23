package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	TypeGithub           = 1
	GithubClientID       = "623037fcf1a6cb4ad6d8"
	GithubClientSecret   = "7ccd7c57dce15c44deee8760f275085afe567708"
	GithubAccessTokenURL = "https://github.com/login/oauth/access_token"
	GithubUserAPI        = "https://api.github.com/user"
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

func (e *User) SignIn(ctx *gin.Context) {
	var payload SignInPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	if payload.Type == TypeGithub {
		tokenResp, err := getGitHubToken(payload.Code)
		if err != nil {
			tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
			return
		}
		githubProfile, err := getGithubProfile(tokenResp.AccessToken)
		if err != nil {
			tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
			return
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
			)
			if err != nil {
				tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
				return
			}
			user = createdUser
		}
		token, err := tools.CreateToken(user.UserID.Hex())
		if err != nil {
			tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
			return
		}
		data := SignInRes{
			User:  user,
			Token: token,
		}
		tools.RespSuccess(ctx, data)
		return
	}
	tools.RespFail(ctx, consts.FailCode, "未知登录类型", nil)
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
