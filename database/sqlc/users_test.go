package db

import (
	"context"
	"moneytor/utils"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

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
