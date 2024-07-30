package suite

import (
	"fmt"
	"testing"

	auth1 "github.com/FelishaK/authProto/gen/auth"
	suite "github.com/FelishaK/grpcAuth/tests"
	"github.com/stretchr/testify/require"
)

func TestGetAuthedUser(t *testing.T) {
	ctx, st := suite.New(t)

	//mock a user
	_, err := st.AuthClient.Register(ctx, &auth1.RegisterRequest{
		FirstName: "Test",
		LastName: "test",
		Password: "123456",
		ConfirmPassword: "123456",
		Email: "mock1@gmail.com",
	})

	if err != nil {
		t.Fatal(err.Error())
	}

	tokenResp, err := st.AuthClient.Login(ctx, &auth1.LoginRequest{
		Email: "mock1@gmail.com",
		Password: "123456",
	})

	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Printf("\n%+v", tokenResp)

	userResp, err := st.AuthClient.GetAuthedUser(ctx, &auth1.TokenRequest{
		Token: tokenResp.AccessToken,
	})

	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Printf("Authed user: %+v", userResp)
	require.NoError(t, err)
}