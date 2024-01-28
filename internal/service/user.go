package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/model"
	"knowledge-service/pkg/consts"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserService struct{}

func (e *UserService) WechatLogin(
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
			consts.TypeWechat,
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

func (e *UserService) GithubLogin(ctx *gin.Context, code string) (model.User, error) {
	tokenResp, err := e.getGitHubToken(code)
	if err != nil {
		return model.User{}, err
	}
	githubProfile, err := e.GetGithubProfile(tokenResp.AccessToken)
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
			consts.TypeGithub,
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

func (e *UserService) getGitHubToken(code string) (entity.GitHubTokenSuccessResp, error) {
	params := entity.GitHubTokenPayload{
		ClientID:     consts.GithubClientID,
		ClientSecret: consts.GithubClientSecret,
		Code:         code,
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return entity.GitHubTokenSuccessResp{}, err
	}
	req, err := http.NewRequest("POST", consts.GithubAccessTokenURL, bytes.NewBuffer(jsonParams))
	if err != nil {
		return entity.GitHubTokenSuccessResp{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return entity.GitHubTokenSuccessResp{}, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return entity.GitHubTokenSuccessResp{}, errors.New("GitHub API request failed")
	}
	var tokenResp entity.GitHubTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return entity.GitHubTokenSuccessResp{}, err
	}
	if len(tokenResp.Error) > 0 {
		return entity.GitHubTokenSuccessResp{}, errors.New(tokenResp.ErrorDescription)
	}
	successResp := entity.GitHubTokenSuccessResp{
		AccessToken: tokenResp.AccessToken,
		Scope:       tokenResp.Scope,
		TokenType:   tokenResp.TokenType,
	}
	return successResp, nil
}

func (e *UserService) GetGithubProfile(token string) (entity.GithubProfileResp, error) {
	req, err := http.NewRequest("GET", consts.GithubUserAPI, nil)
	if err != nil {
		return entity.GithubProfileResp{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return entity.GithubProfileResp{}, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return entity.GithubProfileResp{}, errors.New("GitHub API request failed")
	}
	var profile entity.GithubProfileResp
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return entity.GithubProfileResp{}, err
	}
	return profile, nil
}
