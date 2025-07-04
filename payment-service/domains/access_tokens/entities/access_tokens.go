package entities

import "time"

type AccessToken struct {
	Id        string
	UserId    int
	Revoked   bool
	ExpiresAt time.Time
}
