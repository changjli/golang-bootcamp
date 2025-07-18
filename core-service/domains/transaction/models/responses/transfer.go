package responses

type TransferResponse struct {
	TransactionID int    `json:"transactionID"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}
