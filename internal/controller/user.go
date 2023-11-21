package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"knowledge-service/internal/api"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/model"
	"knowledge-service/internal/service"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (e *UserController) GetProfile(ctx *gin.Context) {
	var query api.GetProfileQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if query.UserID != "" {
		userID = query.UserID
	} else if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	} else {
		tools.RespFail(ctx, consts.Fail, "当前用户不存在", nil)
		return
	}
	userD := dao.UserDAO{}
	user, err := userD.FindByUserID(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, user)
}

func (e *UserController) UpdateProfile(ctx *gin.Context) {
	var payload api.UpdateProfilePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	} else {
		tools.RespFail(ctx, consts.Fail, "当前用户不存在", nil)
		return
	}
	userD := dao.UserDAO{}
	user, err := userD.Update(ctx, userID, payload.Nickname, payload.Avatar)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, user)
}

func (e *UserController) Login(ctx *gin.Context) {
	var payload api.LoginPayload
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
	res := api.LoginRes{
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
	var result api.YDWechatQRCodeResp
	if err = json.Unmarshal(body, &result); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	if result.Code != 0 {
		tools.RespFail(ctx, consts.Fail, result.Msg, nil)
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
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
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
