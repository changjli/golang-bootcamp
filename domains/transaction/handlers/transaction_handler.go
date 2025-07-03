package handlers

import (
	"errors"
	"login-api/domains/transaction"
	"login-api/domains/transaction/models/requests"
	"login-api/domains/users/entities"
	"net/http"
	"strconv"

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

func (h *TransactionHandler) TopUp(ctx *gin.Context) {
	userID, err := h.getAuthenticatedUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req requests.TopUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	response, err := h.transactionUsecase.TopUp(ctx, userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) Transfer(ctx *gin.Context) {
	fromUserID, err := h.getAuthenticatedUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req requests.TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	response, err := h.transactionUsecase.Transfer(ctx, fromUserID, &req)
	if err != nil {
		// Specific error handling for known business logic failures.
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// As per openapi.yaml, 202 Accepted is used for initiated transfers.
	ctx.JSON(http.StatusAccepted, response)
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

func (h *TransactionHandler) GetHistory(ctx *gin.Context) {
	userID, err := h.getAuthenticatedUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Safely parse pagination parameters from the query string with defaults.
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	response, err := h.transactionUsecase.GetHistory(ctx, userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
