package authdto

import (
	"errors"
	"fmt"
	"os"

	auth1 "github.com/FelishaK/authProto/gen/auth"
	jwt "github.com/FelishaK/grpcAuth/internal/lib/jwt"
	jwtvalidate "github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Token string
}

func TokenReqToDTO(rq *auth1.TokenRequest) *Token {
	return &Token{
		Token: rq.Token,
	}
}

func (t *Token) Validate() (*jwtvalidate.Token, error) {
	const op = "authdto.Validate"

	token, err := jwtvalidate.ParseWithClaims(t.Token, &jwt.CustomTokenClaims{}, func(parsedToken *jwtvalidate.Token) (any, error) {
		// reading public key in bytes
		publicKeyBytes, err := os.ReadFile(fmt.Sprintf("%s/public_key.pem", os.Getenv("PEM_FOLDER_PATH")))
	
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		publicKey, err := jwtvalidate.ParseRSAPublicKeyFromPEM(publicKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("%s: %w",op, err)
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return token, nil
	} else if errors.Is(err, jwtvalidate.ErrTokenExpired) {
		return nil, err
	} else {
		return nil, errors.New("invalid token")
	}
}