package requests

type ShortenUrlRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
}
