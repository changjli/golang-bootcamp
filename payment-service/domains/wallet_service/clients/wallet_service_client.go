package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"payment-service/wizards"
	"time"

	"github.com/gin-gonic/gin"
)

// WalletServiceClientImpl is the implementation of WalletServiceClient.
type WalletServiceClientImpl struct {
	client         *http.Client
	coreServiceURL string // e.g., "http://localhost:8080"
}

// NewWalletServiceClient creates a new client for the wallet service.
func NewWalletServiceClient(cfg *wizards.Config) *WalletServiceClientImpl {
	return &WalletServiceClientImpl{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		coreServiceURL: cfg.CoreServiceURL,
	}
}

// VerifyBalanceRequest is the request body for the internal verification endpoint.
type VerifyBalanceRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

// ErrorResponse is the expected error structure from the core-service.
type ErrorResponse struct {
	Error string `json:"error"`
}

// VerifyBalance makes an HTTP call to the core-service to verify a user's balance.
func (c *WalletServiceClientImpl) VerifyBalance(ctx *gin.Context, userID string, amount float64) error {
	// 1. Create the request body
	reqBody := VerifyBalanceRequest{
		UserID: userID,
		Amount: amount,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 2. Create the HTTP request
	url := fmt.Sprintf("%s/api/internal/wallets/verify_balance", c.coreServiceURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call core-service: %w", err)
	}
	defer resp.Body.Close()

	// 4. Handle the response
	if resp.StatusCode == http.StatusOK {
		// Success case
		return nil
	}

	// Error case
	respBody, _ := io.ReadAll(resp.Body)
	var errResp ErrorResponse
	if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != "" {
		return errors.New(errResp.Error)
	}

	return fmt.Errorf("core-service returned status %d", resp.StatusCode)
}
