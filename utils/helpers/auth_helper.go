package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserId      int
	Jti         string
	TokenString string
}

// GetAuthenticatedClaims extracts the user ID from the Gin context's claims.
// It's a shared utility function to be used by any handler that needs the current user's ID.
func GetAuthenticatedClaims(ctx *gin.Context) (*Claims, error) {
	ctxVal, exists := ctx.Get("claims")
	if !exists {
		return nil, errors.New("claims not found in context, user not authenticated")
	}

	claims, ok := ctxVal.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims format in context")
	}

	return claims, nil
}
