package responses

import "time"

type TransactionDetail struct {
	TransactionID string    `json:"transactionID"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	From          string    `json:"from"`
	To            string    `json:"to"`
	Timestamp     time.Time `json:"timestamp"`
	Status        string    `json:"status"`
}

type Pagination struct {
	CurrentPage int   `json:"currentPage"`
	TotalPages  int   `json:"totalPages"`
	TotalItems  int64 `json:"totalItems"`
}

type TransactionHistoryResponse struct {
	Transactions []TransactionDetail `json:"transactions"`
	Pagination   Pagination          `json:"pagination"`
}
