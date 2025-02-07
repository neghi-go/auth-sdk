package users

import (
	"context"
	"os"
	"testing"

	"github.com/neghi-go/database/mongodb"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

var test_url string

func TestMain(m *testing.M) {
	image := testcontainers.ContainerRequest{
		Image:        "mongo:8.0",
		ExposedPorts: []string{"27017/tcp"},
	}

	client, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: image,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	test_url, _ = client.Endpoint(context.Background(), "")

	exit := m.Run()

	_ = testcontainers.TerminateContainer(client)
	os.Exit(exit)
}

func TestUser(t *testing.T) {
	mgd, err := mongodb.New("mongodb://"+test_url, "test-db")
	require.NoError(t, err)
	userModel, err := mongodb.RegisterModel(mgd, "users", UserModel{})
	require.NoError(t, err)
	userStruct := New(userModel)

	t.Run("Create User With Email Only", func(t *testing.T) {
		u, err := userStruct.CreateUser(context.Background(), "jon@doe.com")
		require.NoError(t, err)
		err = userStruct.store.Save(*u)
		require.NoError(t, err)
		require.NotEmpty(t, u)
	})
	t.Run("Create User With Email and Password", func(t *testing.T) {
		u, err := userStruct.CreateUser(context.Background(), "jane@doe.com", SetPassword("Pass1234."))
		require.NoError(t, err)
		err = userStruct.store.Save(*u)
		require.NoError(t, err)
		require.NotEmpty(t, u)
	})

	t.Run("Retrive User that Exists", func(t *testing.T) {
		u, err := userStruct.RetrieveUser(context.Background(), "jon@doe.com")
		require.NoError(t, err)
		require.NotEmpty(t, u)
	})
	t.Run("Retrive User that Does not Exist", func(t *testing.T) {
		u, err := userStruct.RetrieveUser(context.Background(), "julie@doe.com")
		require.Error(t, err)
		require.Empty(t, u)
	})
}
