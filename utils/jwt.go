package utils

import (
	"github.com/dgrijalva/jwt-go"
	"learning/conf"
	"time"
)

var jwtSecret = []byte(conf.AppConfig.JwtSecret)

type Claims struct {
	Id uint `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(id uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(48 * time.Hour)

	claims := Claims{
		id,
		jwt.StandardClaims{
			IssuedAt:  nowTime.Unix(),
			ExpiresAt: expireTime.Unix(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
