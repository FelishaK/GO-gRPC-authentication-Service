package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomTokenClaims struct {
	IsRefresh bool
	jwt.RegisteredClaims
}

func NewTokenClaims(uid string, email string, exp *jwt.NumericDate, isRefresh bool) *CustomTokenClaims{

	if isRefresh {
		return &CustomTokenClaims{
			true,
			jwt.RegisteredClaims{
				ExpiresAt: exp,
				IssuedAt: jwt.NewNumericDate(time.Now()),
				Subject:  uid,
			},
		}
	}
	return &CustomTokenClaims{
			false,
			jwt.RegisteredClaims{
				ExpiresAt: exp,
				IssuedAt: jwt.NewNumericDate(time.Now()),
				Subject: email,
				Issuer: uid,
			},
		}
}

func NewToken(userId, email string, exp time.Duration, isRefresh bool) (string, error) {
	const op = "jwt.NewToken"

	expTime := jwt.NewNumericDate(time.Now().Add(exp))

	// create claims depending on jwt token type
	customClaims := NewTokenClaims(userId, email, expTime, isRefresh)
	fmt.Printf("NewTokenClaims %+v", customClaims)
	
	// read privateKey
	privateKeyBytes, err := os.ReadFile(fmt.Sprintf("%s/private_key.pem",  os.Getenv("PEM_FOLDER_PATH")))
	
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// create a token with claims
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, customClaims)

	return t.SignedString(privateKey)
}