package api

import (
	db "moneytor/database/sqlc"
	"moneytor/token"
	"moneytor/utils"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store      db.Store
	tokenMaker token.Maker
	config     utils.Config
	router     *gin.Engine
}

func NewServer(store db.Store, config utils.Config, tokenMaker token.Maker) *Server {
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.POST("users", server.register)
	v1.POST("users/login", server.login)
	v1.POST("users/refresh", server.refresh)

	auth := v1.Group("/").Use(authMiddleware(tokenMaker))

	auth.POST("accounts", server.createAccount)
	auth.GET("accounts", server.listAccounts)
	auth.POST("categories", server.createCategory)
	auth.GET("categories", server.listCategories)
	auth.GET("transaction-types", server.listTransactionTypes)
	auth.POST("entries", server.createEntry)
	auth.GET("entries", server.listEntries)
	auth.GET("currencies", server.listCurrencies)

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
