package suite

import (
	"fmt"
	"testing"
	"time"

	auth1 "github.com/FelishaK/authProto/gen/auth"
	"github.com/FelishaK/grpcAuth/internal/lib/jwt"
	suite "github.com/FelishaK/grpcAuth/tests"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRefresh(t *testing.T) {
	ctx, st := suite.New(t)

	tEmail := gofakeit.Email()
	userResp, err := st.AuthClient.Register(ctx, &auth1.RegisterRequest{
		FirstName: gofakeit.Name(),
		LastName: gofakeit.LastName(),
		Password: "123456",
		ConfirmPassword: "123456",
		Email: tEmail,
	})

	if err != nil {
		t.Error("Failed to register user")
	}

	token, err := jwt.NewToken(userResp.UserId, tEmail, time.Hour * 1, true)
	fmt.Printf("TOKEN: %s", token)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = st.AuthClient.Refresh(ctx, &auth1.TokenRequest{
		Token: token,
	})

	if err != nil {
		t.Error(err.Error())
	}

	Exptoken, _ := jwt.NewToken(userResp.UserId, tEmail, time.Millisecond * 1, true)

	_, err = st.AuthClient.Refresh(ctx, &auth1.TokenRequest{
		Token: Exptoken,
	})

	if err != nil {
		st, ok := status.FromError(err)

		if !ok {
		t.Errorf("Error is not supposed to be a gRPC status error: %v", err)

		if st.Code() != codes.PermissionDenied {
			t.Errorf("Expected status code %v, got %v", codes.PermissionDenied, st.Code())
		}
	}
	}
}