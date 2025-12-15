package dto

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required"`
}

type LogoutRequest struct {
	UID uint64 `json:"uid" binding:"required"`
}

type UserToken struct {
	UID      uint64 `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
