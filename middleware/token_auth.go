package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkUserAuth(ctx, func(err error) {
			errMsg := "token无效或已过期"
			if err != nil {
				errMsg = err.Error()
			}
			tools.RespFail(ctx, consts.InvalidToken, errMsg, nil)
			ctx.Abort()
		})
	}
}

func UseToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkUserAuth(ctx, func(err error) {
			ctx.Next()
		})
	}
}

type Callback func(err error)

type TokenPayload struct {
	Token string `json:"token"`
}

// token 优先级：payload > query > header
func checkUserAuth(ctx *gin.Context, callback Callback) {
	tokenStr := ctx.GetHeader("Authorization")
	tokenQuery := ctx.Query("token")
	if tokenQuery != "" {
		tokenStr = tokenQuery
	}
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		callback(err)
		return
	}
	if len(body) != 0 {
		var tokenPayload TokenPayload
		if err := json.Unmarshal(body, &tokenPayload); err != nil {
			callback(err)
			return
		}
		// 重新设置请求体内容
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		if tokenPayload.Token != "" {
			tokenStr = tokenPayload.Token
		}
	}
	if tokenStr == "" {
		callback(nil)
		return
	}
	token, err := tools.ParseToken(tokenStr)
	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		newTokenStr, err := getNewTokenStr(tokenStr)
		if err != nil {
			callback(err)
			return
		}
		newToken, err := tools.ParseToken(newTokenStr)
		if err != nil {
			callback(err)
			return
		}
		if !newToken.Valid {
			callback(nil)
			return
		}
		token = newToken
	} else if err != nil || !token.Valid {
		callback(err)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := claims["uid"].(string)
		err := refreshToken(userID, tokenStr, claims)
		if err != nil {
			callback(err)
			return
		}
		ctx.Set("uid", userID)
	}
	ctx.Next()
}

func refreshToken(userID string, originalTokenStr string, claims jwt.MapClaims) error {
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return err
	}
	remainTime := expirationTime.Sub(time.Now())
	if remainTime > time.Hour*24*6 {
		return nil
	}
	rdsInst := tools.Redis{}
	rds := rdsInst.GetRDS()
	newTokenStr, err := tools.CreateToken(userID)
	if err != nil {
		return err
	}
	rds.Set(originalTokenStr, newTokenStr, time.Hour*24*7)
	return nil
}

func getNewTokenStr(originalTokenStr string) (string, error) {
	rdsInst := tools.Redis{}
	rds := rdsInst.GetRDS()
	newTokenStr, err := rds.Get(originalTokenStr).Result()
	if err != nil {
		return "", err
	}
	return newTokenStr, nil
}
