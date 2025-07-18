package responses

type TopUpResponse struct {
	TransactionID int     `json:"transactionID"`
	Message       string  `json:"message"`
	NewBalance    float64 `json:"newBalance"`
}
