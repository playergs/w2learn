package utils

import (
	"errors"
	"time"
	"w2learn/internal/dto"
	"w2learn/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var globalSecret string

func InitJwt(secret string) {
	if globalSecret == "" {
		globalSecret = secret
	} else {
		logger.Error("JWT secret has already been initialized")
	}
}

func GenerateJwtToken(token *dto.UserToken) (string, error) {
	if token == nil {
		return "", errors.New("token is nil")
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString([]byte(globalSecret))
}

func ParseJWT(tokenStr string) (*dto.UserToken, error) {
	token := &dto.UserToken{}
	claims, err := jwt.ParseWithClaims(tokenStr, token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(globalSecret), nil
	})

	if err != nil || !claims.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return token, nil
}

func IsTokenExpired(token *dto.UserToken) bool {
	if token.ExpiresAt.Before(time.Now()) {
		logger.Error("Token expired",
			zap.String("username", token.Username),
			zap.Uint64("user_id", token.UID),
		)
		return true
	}
	return false
}
