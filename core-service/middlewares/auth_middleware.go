package middlewares

import (
	accesstokens "core-service/domains/access_tokens"
	"core-service/domains/users/entities"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	AccessTokenUsecase accesstokens.AccessTokenUsecaseInterface
}

func NewAuthMiddleware(accessTokenUsecase accesstokens.AccessTokenUsecaseInterface) *AuthMiddleware {
	return &AuthMiddleware{
		AccessTokenUsecase: accessTokenUsecase,
	}
}

func (m *AuthMiddleware) HandleProtectedRoutes(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))

	claims := &entities.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("harusnyambildarienvini"), nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// == database method ==
	err = m.AccessTokenUsecase.Validate(c, claims.Jti)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// == cache method ==
	// // Check access token
	// _, found := wizards.Cache.Get(claims.Jti)
	// if found {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	// 	c.Abort()
	// 	return
	// }

	c.Set("claims", claims)

	c.Next()
}
