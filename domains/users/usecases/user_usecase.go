package usecases

import (
	"fmt"
	accesstokens "login-api/domains/access_tokens"
	"login-api/domains/users/entities"
	"login-api/domains/users/models/requests"
	"login-api/domains/users/models/responses"
	"login-api/domains/users/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserUseCase struct {
	UserRepository     repositories.UserRepositoryInterface
	AccessTokenUseCase accesstokens.AccessTokenUsecaseInterface
}

func NewUserUseCase(userRepository repositories.UserRepositoryInterface, accessTokenUseCase accesstokens.AccessTokenUsecaseInterface) *UserUseCase {
	return &UserUseCase{
		UserRepository:     userRepository,
		AccessTokenUseCase: accessTokenUseCase,
	}
}

func (uc *UserUseCase) Login(ctx *gin.Context, request requests.UserLoginRequest) (*responses.UserLoginResponse, error) {
	user, err := uc.UserRepository.FindByUsernameAndPassword(ctx, request.Username, request.Password)
	if err != nil {
		return nil, fmt.Errorf("User not authorized")
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	jti := uuid.New().String()

	claims := &entities.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Username: user.Username,
		Jti:      jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("harusnyambildarienvini"))
	if err != nil {
		return nil, fmt.Errorf("User not authorized")
	}

	response := &responses.UserLoginResponse{
		AccessToken: tokenString,
	}

	// Store access token
	err = uc.AccessTokenUseCase.Create(ctx, jti, user.Id, time.Now().Add(5*time.Minute))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (uc *UserUseCase) Me(ctx *gin.Context) (*responses.UserMeResponse, error) {
	ctxVal, _ := ctx.Get("claims")

	claims := ctxVal.(*entities.Claims)

	user, err := uc.UserRepository.FindByUsername(ctx, claims.Username)
	if err != nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	response := &responses.UserMeResponse{
		Id:       user.Id,
		Username: user.Username,
		Name:     user.Name,
	}

	return response, nil
}

func (uc *UserUseCase) Logout(ctx *gin.Context) error {
	ctxVal, _ := ctx.Get("claims")

	claims := ctxVal.(*entities.Claims)

	// == database method ==
	err := uc.AccessTokenUseCase.Revoke(ctx, claims.Jti)
	if err != nil {
		return err
	}

	// == cache method ==
	// // Revoke access token
	// durationUntilExpiry := time.Until(claims.ExpiresAt.Time)

	// if durationUntilExpiry <= 0 {
	// 	return nil
	// }

	// wizards.Cache.Set(claims.Jti, true, durationUntilExpiry)

	return nil
}
