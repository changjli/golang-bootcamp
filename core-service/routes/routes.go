package routes

import (
	transaction "core-service/domains/transaction/handlers"
	"core-service/domains/users/handlers"
	wallet "core-service/domains/wallet/handlers"
	"core-service/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(userHttp handlers.UserHttpInterface, walletHttp *wallet.WalletHandler, transactionHttp *transaction.TransactionHandler, authMiddleware *middlewares.AuthMiddleware) *gin.Engine {
	route := gin.Default()

	api := route.Group("/api")

	api.POST("/auth/login", userHttp.Login)

	api.Use(authMiddleware.HandleProtectedRoutes)

	api.GET("/auth/me", userHttp.Me)
	api.POST("/auth/logout", userHttp.Logout)

	wallets := api.Group("/wallets")
	{
		wallets.GET("/balance", walletHttp.GetBalance)
	}

	transactions := api.Group("/transactions")
	{
		transactions.POST("/topup", transactionHttp.TopUp)
		transactions.POST("/transfer", transactionHttp.Transfer)
		transactions.GET("/history", transactionHttp.GetHistory)
	}

	return route
}
