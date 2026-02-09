package security

import "github.com/golang-jwt/jwt/v5"

type RequestClaims struct {
	UserID   int `json:"sub"`
	jwt.RegisteredClaims
}