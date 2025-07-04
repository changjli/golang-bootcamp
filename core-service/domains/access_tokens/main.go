package accesstokens

import (
	"core-service/domains/access_tokens/entities"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessTokenRepositoryInterface interface {
	Save(ctx *gin.Context, token *entities.AccessToken) error

	FindByID(ctx *gin.Context, id string) (*entities.AccessToken, error)

	Revoke(ctx *gin.Context, id string) error
}

type AccessTokenUsecaseInterface interface {
	Create(ctx *gin.Context, tokenID string, userID int, expiresAt time.Time) error
	Validate(ctx *gin.Context, tokenID string) error
	Revoke(ctx *gin.Context, tokenID string) error
}
