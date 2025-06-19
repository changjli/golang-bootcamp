package routes

import (
	"login-api/domains/users/handlers"
	"login-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(userHttp handlers.UserHttpInterface, authMiddleware *middlewares.AuthMiddleware) *gin.Engine {
	route := gin.Default()

	api := route.Group("/api")

	api.POST("/auth/login", userHttp.Login)

	api.Use(authMiddleware.HandleProtectedRoutes)

	api.GET("/auth/me", userHttp.Me)
	api.POST("/auth/logout", userHttp.Logout)

	return route
}
