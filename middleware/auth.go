package middleware

import (
	"currency-converter/security"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// type userCtxKeyType struct{}

// var userCtxKey = userCtxKeyType{}

type AuthMiddleware struct {
	tokenSvc *security.TokenService
}

func NewAuthMiddleware(tokenSvc *security.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenSvc: tokenSvc,
	}
}

func (a *AuthMiddleware) Handle() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {

		token, err := a.extractBearerToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		claims, err := a.tokenSvc.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// Extra safety checks (optional but recommended)
		if claims.UserID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user is unauhorised",
			})
			c.Abort()
			return
		}
		// c.Set(userCtxKey, claims)

		c.Next()
	})
}

func (*AuthMiddleware) extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header")
	}

	return parts[1], nil
}
