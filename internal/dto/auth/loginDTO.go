package authdto

import (
	"errors"
	"fmt"

	auth1 "github.com/FelishaK/authProto/gen/auth"
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

func LoginReqToDTO(rq *auth1.LoginRequest) *LoginRequest {
	return &LoginRequest{
		Email: rq.GetEmail(),
		Password: rq.GetPassword(),
	}
}

func (lr *LoginRequest) Validate() error {
	if lr.Email == "" {
		fmt.Println("HERE1")
		return errors.New("email is required")
	}
	if lr.Password == "" {
		fmt.Println("HERE2")
		return errors.New("password id required")
	}
	return nil
}
