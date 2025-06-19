package usecases

import (
	"errors"
	"fmt"
	accesstokens "login-api/domains/access_tokens"
	"login-api/domains/access_tokens/entities"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrNotFound = errors.New("record not found")
var ErrTokenInvalid = errors.New("token is invalid or does not exist")
var ErrTokenRevoked = errors.New("token has been revoked")

type AccessTokenUsecase struct {
	accessTokenRepo accesstokens.AccessTokenRepositoryInterface
}

func NewAccessTokenUsecase(repo accesstokens.AccessTokenRepositoryInterface) *AccessTokenUsecase {
	return &AccessTokenUsecase{
		accessTokenRepo: repo,
	}
}

func (u *AccessTokenUsecase) Create(ctx *gin.Context, tokenID string, userID int, expiresAt time.Time) error {
	accessToken := &entities.AccessToken{
		Id:        tokenID,
		UserId:    userID,
		Revoked:   false,
		ExpiresAt: expiresAt,
	}

	if err := u.accessTokenRepo.Save(ctx, accessToken); err != nil {
		return fmt.Errorf("failed to save token to repository: %w", err)
	}

	return nil
}

func (u *AccessTokenUsecase) Validate(ctx *gin.Context, tokenID string) error {
	repoToken, err := u.accessTokenRepo.FindByID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrTokenInvalid
		}
		return fmt.Errorf("repository check failed: %w", err)
	}

	if repoToken.Revoked {
		return ErrTokenRevoked
	}

	return nil
}

func (u *AccessTokenUsecase) Revoke(ctx *gin.Context, tokenID string) error {
	return u.accessTokenRepo.Revoke(ctx, tokenID)
}
