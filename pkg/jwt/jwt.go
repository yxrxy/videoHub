package jwt

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/pkg/constants"
)

var (
	Secret []byte
	once   sync.Once
)

func InitSecret() {
	once.Do(func() {
		if config.JWT == nil {
			panic("JWT config is not initialized")
		}
		Secret = []byte(config.JWT.SecretKey)
	})
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64) (string, error) {
	InitSecret()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(constants.TokenExpiry).Unix(),
	})
	return token.SignedString(Secret)
}

func GenerateRefreshToken(userID int64) (string, error) {
	InitSecret()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	return token.SignedString(Secret)
}

func ParseRefreshToken(tokenString string) (*Claims, error) {
	InitSecret()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func ParseToken(tokenString string) (*Claims, error) {
	InitSecret()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
