package routes

import (
	"login-api/domains/users/handlers"
	"login-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(userHttp handlers.UserHttpInterface) *gin.Engine {
	route := gin.Default()

	api := route.Group("/api")

	api.POST("/auth/login", userHttp.Login)

	api.Use(middlewares.HandleProtectedRoutes)

	api.GET("/auth/me", userHttp.Me)

	return route
}
