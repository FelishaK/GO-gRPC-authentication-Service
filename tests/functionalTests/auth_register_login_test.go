package suite

import (
	"fmt"
	"log"
	"testing"

	auth1 "github.com/FelishaK/authProto/gen/auth"
	suite "github.com/FelishaK/grpcAuth/tests"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Gender struct {
	Gender string
}

func TestRegisterAndLogin(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	firstName := gofakeit.Name()
	lastName := gofakeit.LastName()

	_, err := st.AuthClient.Register(ctx, &auth1.RegisterRequest{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		Password:        "12345678",
		ConfirmPassword: "12345678",
		Gender:          auth1.Gender_Female.Enum(),
	})

	if err != nil {
		t.Errorf(err.Error())
	}

	logResp, err := st.AuthClient.Login(ctx, &auth1.LoginRequest{
		Email:    email,
		Password: "12345678",
	})
	
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Printf("LOGIN RESPONSE: %+v\n", logResp)
	require.NoError(t, err)
}

func TestRegister(t *testing.T) {
	ctx, st := suite.New(t)
	var successTests = []struct {
		FirstName       string
		LastName        string
		Email           string
		Password        string
		ConfirmPassword string
		Gender          string
	}{
		{gofakeit.FirstName(), gofakeit.LastName(), gofakeit.Email(), "123456", "123456", auth1.Gender_Female.Enum().String()},
		{gofakeit.FirstName(), gofakeit.LastName(), gofakeit.Email(), "123456", "123456", ""},
		{"Test", "Test", "test@gmail.com", "123456", "123456", ""},
	}

	var FailTests = []struct {
		FirstName       string
		LastName        string
		Email           string
		Password        string
		ConfirmPassword string
		Gender          string
	}{
		{"", gofakeit.LastName(), gofakeit.Email(), "123456", "123456", auth1.Gender_Female.Enum().String()},
		{gofakeit.Name(), "", gofakeit.Email(), "123456", "123456", auth1.Gender_Male.Enum().String()},
		{gofakeit.Name(), gofakeit.LastName(), "", "123456", "123456", auth1.Gender_Male.Enum().String()},
		{gofakeit.Name(), gofakeit.LastName(), gofakeit.Email(), "", "123456", auth1.Gender_Male.Enum().String()},
		{gofakeit.Name(), gofakeit.LastName(), gofakeit.Email(), "123456", "", auth1.Gender_Male.Enum().String()},
		{gofakeit.Name(), gofakeit.LastName(), gofakeit.Email(), "123456", "12333333", auth1.Gender_Male.Enum().String()},
	}

	for _, tt := range successTests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			_, err := st.AuthClient.Register(ctx, &auth1.RegisterRequest{
				FirstName:       tt.FirstName,
				LastName:        tt.LastName,
				Email:           tt.Email,
				Password:        tt.Password,
				ConfirmPassword: tt.ConfirmPassword,
				Gender:          auth1.Gender_Female.Enum(),
			})

			if err != nil {
				t.Errorf("%s", err.Error())
			}
		})
	}

	for _, tt := range FailTests {
		t.Run("", func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &auth1.RegisterRequest{
				FirstName:       tt.FirstName,
				LastName:        tt.LastName,
				Email:           tt.Email,
				Password:        tt.Password,
				ConfirmPassword: tt.ConfirmPassword,
				Gender:          auth1.Gender_Female.Enum(),
			})

			if err == nil {
				t.Fatal("Expected error, but got nil")
			}
			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("Error is not supposed to be a gRPC status error: %v", err)
			}
			if st.Code() != codes.InvalidArgument {
				t.Errorf("Expected status code %v, got %v", codes.InvalidArgument, st.Code())
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctx, st := suite.New(t)

	var failTests = []struct {
		Email    string
		Password string
	}{
		{gofakeit.Email(), ""},
		{"", "123456"},
	}

	for _, tc := range failTests {
		t.Run("", func(t *testing.T) {
			_, err := st.AuthClient.Login(ctx, &auth1.LoginRequest{
				Email:    tc.Email,
				Password: tc.Password,
			})

			if err == nil {
				log.Fatal("Expected error, but got nil")
			}

			st, ok := status.FromError(err)
			if !ok {
				log.Fatalf("Error is not supposed to be a gRPC status error: %v", err)
			}
			if st.Code() != codes.InvalidArgument {
				t.Errorf("Expected status code %v, got %v", codes.InvalidArgument, st.Code())
			}
		})
	}

	_, err := st.AuthClient.Login(ctx, &auth1.LoginRequest{
		Email:    "fsdfdsfd@gmail.com",
		Password: "123456",
	})

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Fatalf("Error is not supposed to be a gRPC status error: %v", err)
		}

		if st.Code() != codes.NotFound {
			t.Errorf("Expected status code %v, got %v", codes.NotFound, st.Code())
		}
	}
}
