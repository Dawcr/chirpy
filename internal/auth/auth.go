package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return tkn.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	tkn, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil || !tkn.Valid {
		return uuid.Nil, errors.New("token is invalid or has expired")
	}

	id, err := tkn.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, errors.New("unable to access users id from claims")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New("error parsing uuid")
	}

	return uid, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("no authorization header")
	}

	if !strings.HasPrefix(tokenString, "Bearer") {
		return "", errors.New("missing bearer prefix")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer")
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return "", errors.New("bearer token is empty")
	}

	return tokenString, nil
}
