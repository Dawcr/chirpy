package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
	t.Run("valid header should return bearer token", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer abc123")
		token, err := GetBearerToken(headers)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token != "abc123" {
			t.Fatalf("expected token 'abc123', got '%v'", token)
		}
	})

	t.Run("missing authorization header should return error", func(t *testing.T) {
		headers := http.Header{}
		_, err := GetBearerToken(headers)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty token after Bearer should return error", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer  ")
		_, err := GetBearerToken(headers)
		if err == nil {
			t.Fatal("expected error for empty token")
		}
	})

	t.Run("header without Bearer prefix should return error", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "NotBearer abc123")
		_, err := GetBearerToken(headers)
		if err == nil {
			t.Fatal("expected error for missing Bearer prefix")
		}
	})

	t.Run("boot.dev tests", func(t *testing.T) {
		tests := []struct {
			name      string
			headers   http.Header
			wantToken string
			wantErr   bool
		}{
			{
				name: "Valid Bearer token",
				headers: http.Header{
					"Authorization": []string{"Bearer valid_token"},
				},
				wantToken: "valid_token",
				wantErr:   false,
			},
			{
				name:      "Missing Authorization header",
				headers:   http.Header{},
				wantToken: "",
				wantErr:   true,
			},
			{
				name: "Malformed Authorization header",
				headers: http.Header{
					"Authorization": []string{"InvalidBearer token"},
				},
				wantToken: "",
				wantErr:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gotToken, err := GetBearerToken(tt.headers)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotToken != tt.wantToken {
					t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
				}
			})
		}
	})
}

func TestValidateJWT(t *testing.T) {
	const testingSecret = "vgiAfdi/lRsoH2ZqNsbt2FYg1j0u8M2u1EFAnHHCtqsb985DZvGZRzXqGeBS3FmbrhmvocPTAJygA4i0wrrMuw=="

	t.Run("valid JWT should return correct userID", func(t *testing.T) {
		userID := uuid.New()
		token, err := MakeJWT(userID, testingSecret, time.Hour)
		if err != nil {
			t.Fatalf("error creating JWT: %s", err)
		}

		gotID, err := ValidateJWT(token, testingSecret)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if gotID != userID {
			t.Errorf("expected user ID %v, got %v", userID, gotID)
		}
	})

	t.Run("expired JWT should return error", func(t *testing.T) {
		userID := uuid.New()
		token, err := MakeJWT(userID, testingSecret, time.Nanosecond)
		if err != nil {
			t.Fatalf("error creating JWT: %s", err)
		}
		time.Sleep(time.Millisecond)

		_, err = ValidateJWT(token, testingSecret)
		if err == nil {
			t.Fatalf("expected expired token warning")
		}
	})

	t.Run("JWT signed with wrong secret should return error", func(t *testing.T) {
		userID := uuid.New()
		altSecret := "qUFtedLvBJk4yIeig/8DPrbAkzPK8JOrd47qc9bW7BumKAgCfj+7hXi09KUcQOTKUymha8Bh9RP7miLh6lbUlw=="
		token, err := MakeJWT(userID, altSecret, time.Nanosecond)
		if err != nil {
			t.Fatalf("error creating JWT: %s", err)
		}

		_, err = ValidateJWT(token, testingSecret)
		if err == nil {
			t.Fatalf("expected wrong signing secret warning")
		}
	})

	t.Run("malformed JWT should return error", func(t *testing.T) {
		userID := uuid.New()
		token, err := MakeJWT(userID, testingSecret, time.Nanosecond)
		if err != nil {
			t.Fatalf("error creating JWT: %s", err)
		}

		_, err = ValidateJWT(token[5:], testingSecret)
		if err == nil {
			t.Fatalf("expected invalid warning for malformed token")
		}
	})

	t.Run("boot.dev tests", func(t *testing.T) {
		userID := uuid.New()
		validToken, _ := MakeJWT(userID, "secret", time.Hour)

		tests := []struct {
			name        string
			tokenString string
			tokenSecret string
			wantUserID  uuid.UUID
			wantErr     bool
		}{
			{
				name:        "Valid token",
				tokenString: validToken,
				tokenSecret: "secret",
				wantUserID:  userID,
				wantErr:     false,
			},
			{
				name:        "Invalid token",
				tokenString: "invalid.token.string",
				tokenSecret: "secret",
				wantUserID:  uuid.Nil,
				wantErr:     true,
			},
			{
				name:        "Wrong secret",
				tokenString: validToken,
				tokenSecret: "wrong_secret",
				wantUserID:  uuid.Nil,
				wantErr:     true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotUserID != tt.wantUserID {
					t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
				}
			})
		}
	})
}

// boot.dev tests
func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
