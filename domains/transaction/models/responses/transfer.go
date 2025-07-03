package responses

type TransferResponse struct {
	TransactionID string `json:"transactionID"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}
