package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPassword(password string) string {
	cost := 10

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	return string(hashedPassword)
}
