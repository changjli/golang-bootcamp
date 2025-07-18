package models

type RegisterResp struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type MeResp struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type LoginResp struct {
	AccessToken string `json:"access_token"`
}
