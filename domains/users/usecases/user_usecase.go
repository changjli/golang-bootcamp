package usecases

import (
	"fmt"
	"login-api/domains/users/entities"
	"login-api/domains/users/models/requests"
	"login-api/domains/users/models/responses"
	"login-api/domains/users/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserUseCase struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUserUseCase(userRepository repositories.UserRepositoryInterface) *UserUseCase {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

func (uc *UserUseCase) Login(ctx *gin.Context, request requests.UserLoginRequest) (*responses.UserLoginResponse, error) {
	user, err := uc.UserRepository.FindByUsernameAndPassword(ctx, request.Username, request.Password)
	if err != nil {
		return nil, fmt.Errorf("User not authorized")
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &entities.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("harusnyambildarienvini"))
	if err != nil {
		return nil, fmt.Errorf("User not authorized")
	}

	response := &responses.UserLoginResponse{
		AccessToken: tokenString,
	}

	return response, nil
}

func (uc *UserUseCase) Me(ctx *gin.Context) (*responses.UserMeResponse, error) {
	value, exist := ctx.Get("username")

	if !exist {
		return nil, fmt.Errorf("Unauthorized")
	}

	username, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}

	user, err := uc.UserRepository.FindByUsername(ctx, string(username))
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
