package utils

import (
	"learning/conf"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(conf.AppConfig.JwtSecret)

type Claims struct {
	Id   uint `json:"id"`
	Role uint `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(id uint, role uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(48 * time.Hour)

	claims := Claims{
		id,
		role,
		jwt.StandardClaims{
			IssuedAt:  nowTime.Unix(),
			ExpiresAt: expireTime.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parsing token
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
