package responses

type ShortenUrlResponse struct {
	ShortUrl  string `json:"short_url"`
	ShortCode string `json:"short_code"`
}
