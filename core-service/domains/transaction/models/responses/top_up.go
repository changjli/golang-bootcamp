package responses

type TopUpResponse struct {
	TransactionID string  `json:"transactionID"`
	Message       string  `json:"message"`
	NewBalance    float64 `json:"newBalance"`
}
