package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

const (
	chosenCost = 12
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), chosenCost)
	if err != nil {
		log.Fatalf("Error hashing password: %s", err)
	}

	return string(hashed), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
