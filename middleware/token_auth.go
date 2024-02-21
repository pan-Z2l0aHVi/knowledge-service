package middleware

import (
	"errors"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkUserAuth(ctx, func() {
			tools.RespFail(ctx, consts.InvalidToken, "token无效或已过期", nil)
			ctx.Abort()
		})
	}
}

func UseToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkUserAuth(ctx, ctx.Next)
	}
}

type Callback func()

func checkUserAuth(ctx *gin.Context, callback Callback) {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		callback()
		return
	}
	token, err := tools.ParseToken(tokenStr)
	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		newTokenStr, err := getNewTokenStr(tokenStr)
		if err != nil {
			callback()
			return
		}
		newToken, err := tools.ParseToken(newTokenStr)
		if err != nil {
			callback()
			return
		}
		if !newToken.Valid {
			callback()
			return
		}
		token = newToken
	} else if err != nil || !token.Valid {
		callback()
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := claims["uid"].(string)
		err := refreshToken(userID, tokenStr)
		if err != nil {
			callback()
			return
		}
		ctx.Set("uid", userID)
	}
	ctx.Next()
}

func refreshToken(userID string, originalTokenStr string) error {
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
