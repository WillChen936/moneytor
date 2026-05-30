package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	RandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := RandomUser(t)

	userGet, err := testStore.GetUser(context.Background(), user.ID)

	require.NoError(t, err)
	require.Equal(t, user.ID, userGet.ID)
	require.Equal(t, user.Username, userGet.Username)
	require.Equal(t, user.Email, userGet.Email)
	require.Equal(t, user.HashedPassword, userGet.HashedPassword)
}

func TestGetUserByEmail(t *testing.T) {
	user := RandomUser(t)

	userGet, err := testStore.GetUserByEmail(context.Background(), user.Email)

	require.NoError(t, err)
	require.Equal(t, user.ID, userGet.ID)
	require.Equal(t, user.Username, userGet.Username)
	require.Equal(t, user.Email, userGet.Email)
}

func RandomUser(t *testing.T) User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(utils.RandomString(8)), bcrypt.DefaultCost)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomString(6),
		Email:          utils.RandomString(6) + "@test.com",
		HashedPassword: string(hashedPassword),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, user.ID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)

	return user
}
