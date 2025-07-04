package helpers

import (
	"errors"
	"payment-service/domains/users/entities"

	"github.com/gin-gonic/gin"
)

// GetAuthenticatedUserID extracts the user ID from the Gin context's claims.
// It's a shared utility function to be used by any handler that needs the current user's ID.
func GetAuthenticatedUserID(ctx *gin.Context) (string, error) {
	ctxVal, exists := ctx.Get("claims")
	if !exists {
		return "", errors.New("claims not found in context, user not authenticated")
	}

	claims, ok := ctxVal.(*entities.Claims)
	if !ok {
		return "", errors.New("invalid claims format in context")
	}

	if string(claims.UserId) == "" {
		return "", errors.New("user ID is empty in claims")
	}

	return string(claims.UserId), nil
}
