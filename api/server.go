package api

import (
	db "moneytor/database/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()
	v1Routes := router.Group("/api/v1")

	v1Routes.GET("health", server.health)
	v1Routes.POST("accounts", server.createAccount)
	v1Routes.GET("accounts", server.listAccounts)
	v1Routes.POST("categories", server.createCategory)
	v1Routes.GET("categories", server.listCategories)
	v1Routes.GET("transaction-types", server.listTransactionTypes)
	v1Routes.POST("entries", server.createEntry)
	v1Routes.GET("entries", server.listEntries)
	v1Routes.GET("currencies", server.listCurrencies)

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func (s *Server) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
