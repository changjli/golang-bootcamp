package routes

import (
	transaction "login-api/domains/transaction/handlers"
	"login-api/domains/users/handlers"
	wallet "login-api/domains/wallet/handlers"
	"login-api/middlewares"

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
		transactions.POST("/pay", transactionHttp.Pay)
		transactions.GET("/history", transactionHttp.GetHistory)
	}

	return route
}
