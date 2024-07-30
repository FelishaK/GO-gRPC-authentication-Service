package suite

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	auth1 "github.com/FelishaK/authProto/gen/auth"
	"github.com/FelishaK/grpcAuth/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient auth1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Duration(time.Minute* 10))

	t.Cleanup(func() {
		//clears db after all tests have passed
		clearDB(cfg, "users")
		t.Helper()
		cancelCtx()
	})

	// creating a new client
	con, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Grpc.GrpcHost, strconv.Itoa(cfg.Grpc.Port)), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:         t,
		Cfg:        cfg,
		AuthClient: auth1.NewAuthClient(con),
	}
}


// deletes collection
func clearDB(config *config.Config, collName string) {
	conStr := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Database.MongoUser, config.Database.MongoPassword, config.Database.MongoHostname, strconv.Itoa(config.Database.MongoPort))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conStr))
	if err != nil {
		panic(err)
	}
	err = client.Database(config.Database.MongoDBName).Collection(collName).Drop(ctx)
	if err != nil {
		panic(err)
	}
	defer cancel()
}

