package requests

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required,min=1"`
}
