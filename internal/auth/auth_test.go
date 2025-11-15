package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJwt(t *testing.T) {
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
}

func TestGetBearerToken(t *testing.T) {

	tests := []struct {
		name      string
		headers   http.Header
		wantErr   bool
		wantToken string
	}{
		{
			name:      "No authorization header",
			headers:   http.Header{},
			wantErr:   true,
			wantToken: "",
		},
		{
			name: "No bearer token",
			headers: http.Header{
				"Authorization": []string{"Token jwt_token"},
			},
			wantErr:   true,
			wantToken: "",
		},
		{
			name: "empty bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			wantErr:   true,
			wantToken: "",
		},
		{
			name: "valid bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer jwt_token"},
			},
			wantErr:   false,
			wantToken: "jwt_token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if err == nil && res != tt.wantToken {
				t.Errorf("GetBearerToken() value = %v, wantValue = %v", res, tt.wantToken)
			}
		})
	}
}
