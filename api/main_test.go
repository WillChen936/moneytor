package api

import (
	db "moneytor/database/sqlc"
	"moneytor/utils"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSecretKey:       utils.RandomString(32),
		AccessTokenDuration:  time.Minute,
		RefreshTokenDuration: time.Hour,
	}

	server, err := NewServer(store, config)

	require.NoError(t, err)

	return server
}

func addAuthorization(t *testing.T, request *http.Request, server *Server, userID int64) {
	tokenStr, _, err := server.tokenMaker.CreateToken(userID, time.Minute)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+tokenStr)
}
