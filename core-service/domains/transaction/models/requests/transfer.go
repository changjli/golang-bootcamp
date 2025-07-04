package requests

type TransferRequest struct {
	ToUserID string  `json:"toUserID" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,min=1"`
}
