package usecases

import (
	accesstokens "core-service/domains/access_tokens"
	"core-service/domains/users"
	"core-service/domains/users/entities"
	"core-service/domains/users/models"
	"fmt"
	"time"
	"utils/helpers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCaseImpl struct {
	UserRepository     users.UserRepository
	AccessTokenUseCase accesstokens.AccessTokenUsecaseInterface
}

func NewUserUseCase(userRepository users.UserRepository, accessTokenUseCase accesstokens.AccessTokenUsecaseInterface) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		UserRepository:     userRepository,
		AccessTokenUseCase: accessTokenUseCase,
	}
}

func (uc *UserUseCaseImpl) Register(ctx *gin.Context, request models.RegisterReq) (*models.RegisterResp, error) {
	// Check if the name is already taken to avoid duplicates.
	_, err := uc.UserRepository.FindbyName(ctx, request.Name)
	if err == nil {
		return nil, fmt.Errorf("name '%s' is already taken", request.Name)
	}

	// Hash the user's password for secure storage.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create a new user entity.
	user := &entities.User{
		Name:         request.Name,
		PasswordHash: string(hashedPassword),
	}

	// Save the new user to the repository.
	savedUser, err := uc.UserRepository.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	// Prepare the response, excluding sensitive data like the password.
	response := &models.RegisterResp{
		UserID: savedUser.ID,
		Name:   savedUser.Name,
		Email:  savedUser.Email,
	}

	return response, nil
}

func (uc *UserUseCaseImpl) Login(ctx *gin.Context, request models.LoginReq) (*models.LoginResp, error) {
	user, err := uc.UserRepository.FindbyEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	jti := uuid.New().String()

	claims := &helpers.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		UserId: user.ID,
		Jti:    jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("harusnyambildarienvini"))
	if err != nil {
		return nil, fmt.Errorf("User not authorized")
	}

	response := &models.LoginResp{
		AccessToken: tokenString,
	}

	// Store access token
	err = uc.AccessTokenUseCase.Create(ctx, jti, user.ID, time.Now().Add(5*time.Minute))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (uc *UserUseCaseImpl) Me(ctx *gin.Context) (*models.MeResp, error) {
	claims, err := helpers.GetAuthenticatedClaims(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.UserRepository.FindbyID(ctx, claims.UserId)
	if err != nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	response := &models.MeResp{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}

	return response, nil
}

func (uc *UserUseCaseImpl) Logout(ctx *gin.Context) error {
	claims, err := helpers.GetAuthenticatedClaims(ctx)
	if err != nil {
		return err
	}

	// == database method ==
	err = uc.AccessTokenUseCase.Revoke(ctx, claims.Jti)
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
