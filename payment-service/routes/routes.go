package routes

import (
	transaction "payment-service/domains/transaction/handlers"
	"payment-service/domains/users/handlers"
	"payment-service/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(userHttp handlers.UserHttpInterface, transactionHttp *transaction.TransactionHandler, authMiddleware *middlewares.AuthMiddleware) *gin.Engine {
	route := gin.Default()

	api := route.Group("/api")

	api.Use(authMiddleware.HandleProtectedRoutes)

	transactions := api.Group("/transactions")
	{
		transactions.POST("/pay", transactionHttp.Pay)
	}

	return route
}
