package middlewares

import (
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			tools.RespFail(ctx, consts.InvalidToken, "token无效或已过期", nil)
			ctx.Abort()
			return
		}
		token, err := tools.ParseToken(tokenStr)
		if err != nil {
			tools.RespFail(ctx, consts.InvalidToken, "token无效或已过期", nil)
			ctx.Abort()
			return
		}
		if !token.Valid {
			tools.RespFail(ctx, consts.InvalidToken, "token无效或已过期", nil)
			ctx.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("uid", claims["uid"])
		}
		ctx.Next()
	}
}
