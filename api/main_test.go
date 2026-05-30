package api

import (
	db "moneytor/database/sqlc"
	"moneytor/token"
	"moneytor/utils"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) (*Server, token.Maker) {
	config := utils.Config{
		TokenSecretKey:       utils.RandomString(32),
		AccessTokenDuration:  time.Minute,
		RefreshTokenDuration: time.Hour,
	}
	tokenMaker, err := token.NewJWTMaker(config.TokenSecretKey)
	require.NoError(t, err)
	return NewServer(store, config, tokenMaker), tokenMaker
}

func addAuthorization(t *testing.T, request *http.Request, maker token.Maker, userID int64) {
	tokenStr, _, err := maker.CreateToken(userID, time.Minute)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+tokenStr)
}
