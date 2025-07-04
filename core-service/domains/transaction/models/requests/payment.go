package requests

type PayRequest struct {
	MerchantID  string  `json:"merchantID" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,min=1"`
	Description string  `json:"description"`
}
