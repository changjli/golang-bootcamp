package responses

type WalletBalanceResponse struct {
	UserID  string  `json:"userID"`  // The ID of the user who owns the wallet
	Balance float64 `json:"balance"` // The current balance of the wallet
}
