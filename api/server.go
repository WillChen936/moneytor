package api

import (
	"fmt"
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

func NewServer(store db.Store, config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

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
	v1.GET("transaction-types", server.listTransactionTypes)
	v1.GET("currencies", server.listCurrencies)

	auth := v1.Group("/").Use(authMiddleware(tokenMaker))

	auth.POST("accounts", server.createAccount)
	auth.GET("accounts", server.listAccounts)
	auth.POST("categories", server.createCategory)
	auth.GET("categories", server.listCategories)
	auth.POST("entries", server.createEntry)
	auth.GET("entries", server.listEntries)

	server.router = router
	return server, nil
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
