package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	user := RandomUser(t)

	arg := CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	session, err := testStore.CreateSession(context.Background(), arg)

	require.NoError(t, err)
	require.True(t, session.ID.Valid)
	require.Equal(t, arg.UserID, session.UserID)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.WithinDuration(t, arg.ExpiresAt, session.ExpiresAt, time.Second)
}

func TestGetSession(t *testing.T) {
	user := RandomUser(t)

	arg := CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	session, err := testStore.CreateSession(context.Background(), arg)
	require.NoError(t, err)

	sessionGet, err := testStore.GetSession(context.Background(), session.ID)

	require.NoError(t, err)
	require.Equal(t, session.ID, sessionGet.ID)
	require.Equal(t, session.UserID, sessionGet.UserID)
	require.Equal(t, session.RefreshToken, sessionGet.RefreshToken)
	require.WithinDuration(t, session.ExpiresAt, sessionGet.ExpiresAt, time.Second)
}

func TestDeleteSession(t *testing.T) {
	user := RandomUser(t)

	arg := CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	session, err := testStore.CreateSession(context.Background(), arg)
	require.NoError(t, err)

	err = testStore.DeleteSession(context.Background(), session.ID)
	require.NoError(t, err)

	sessionGet, err := testStore.GetSession(context.Background(), session.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, sessionGet)
}
