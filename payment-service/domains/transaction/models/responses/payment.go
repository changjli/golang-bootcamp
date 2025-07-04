package responses

type PayResponse struct {
	TransactionID string  `json:"transactionID"`
	Message       string  `json:"message"`
	NewBalance    float64 `json:"newBalance"`
}
