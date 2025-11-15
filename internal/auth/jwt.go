package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// create JWT
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	res, err := token.SignedString(signingKey)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		return "", err
	}

	return res, err
}

// validate JWT
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		log.Printf("Error validating jwt: %s", err)
		return uuid.Nil, err
	}

	idStr, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Error getting id using GetSubject() while validating jwt: %s", err)
		return uuid.Nil, err
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Error parsing uuid while validating jwt: %s", err)
		return uuid.Nil, err
	}
	return id, nil
}
