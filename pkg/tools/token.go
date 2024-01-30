package tools

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTTokenSecretKey = "jwt"
)

func CreateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(time.Hour * 24 * 14).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTTokenSecretKey))
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTTokenSecretKey), nil
	})
}
