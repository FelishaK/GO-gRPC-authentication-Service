package authdto

import (
	"errors"
	"fmt"
	"net/mail"
	"reflect"

	auth1 "github.com/FelishaK/authProto/gen/auth"
)

type RegisterRequest struct {
	FirstName string
	LastName  string
	Password string
	ConfirmPassword string
	Email     string
	Gender    string
}

type RegisterResponse struct {
	userId string
}

func ToDTO(rq *auth1.RegisterRequest) *RegisterRequest {
		vals := reflect.ValueOf(rq).Elem()
		fieldVal := vals.FieldByName("Gender")

		// check if Gender field is inside of the request struct
		if !fieldVal.IsNil() {
			return &RegisterRequest {
					FirstName: rq.GetFirstName(),
					LastName:  rq.GetLastName(),
					Password: rq.GetPassword(),
					ConfirmPassword: rq.GetConfirmPassword(),
					Email:     rq.GetEmail(),
					Gender: rq.GetGender().String(),
			} 
		}
		return &RegisterRequest {
				FirstName: rq.GetFirstName(),
				LastName:  rq.GetLastName(),
				Password: rq.GetPassword(),
				ConfirmPassword: rq.GetConfirmPassword(),
				Email:     rq.GetEmail(),
			}
}

func (rq *RegisterRequest) Validate() error {
	fmt.Printf("%+v", rq)
	if rq.Password != rq.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	if rq.FirstName == "" {
		return errors.New("firstName is required")
	} 
	if rq.LastName == "" {
		return errors.New("lastName is required")
	}
	if rq.Password == "" {
		return errors.New("password is required")
	}
	_, err := mail.ParseAddress(rq.Email)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}