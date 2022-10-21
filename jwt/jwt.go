package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const ExpiresSecond = 60 * 60 * 12

type CustomClaims struct {
	PrimaryKey string `json:"primaryKey,omitempty"`
	jwt.StandardClaims
}

func GenerateToken(id string) string {
	expiresTime := time.Now().Unix() + int64(ExpiresSecond)
	claims := CustomClaims{
		PrimaryKey: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresTime,
			IssuedAt:  time.Now().Unix() - int64(60*10),
		},
	}
	var jwtSecret = []byte(os.Getenv("jwt.secret"))
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		return token
	} else {
		return ""
	}
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TokenSecret")), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.PrimaryKey, nil
	} else {
		return "", errors.New("token valid error")
	}
}
