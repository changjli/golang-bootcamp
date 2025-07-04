package handlers

import (
	"errors"
	"net/http"
	"payment-service/domains/transaction"
	"payment-service/domains/transaction/models/requests"
	"payment-service/domains/users/entities"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUsecase transaction.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase transaction.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{transactionUsecase: transactionUsecase}
}

func (h *TransactionHandler) getAuthenticatedUserID(ctx *gin.Context) (string, error) {
	ctxVal, exists := ctx.Get("claims")
	if !exists {
		return "", errors.New("user not authenticated")
	}

	claims, ok := ctxVal.(*entities.Claims)

	if !ok {
		return "", errors.New("invalid user ID format in token")
	}
	return string(claims.UserId), nil
}

func (h *TransactionHandler) Pay(ctx *gin.Context) {
	userID, err := h.getAuthenticatedUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req requests.PayRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	response, err := h.transactionUsecase.InitiatePayment(ctx, userID, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 200 OK indicates the payment session was successfully created.
	ctx.JSON(http.StatusOK, response)
}
