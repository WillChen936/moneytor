package api

import (
	db "moneytor/database/sqlc"
	"moneytor/token"
	"moneytor/utils"
	"net/http"

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

	allowedOrigins := []string{
		"http://localhost:5500",
		"http://127.0.0.1:5500",
	}
	router.Use(corsMiddleware(allowedOrigins))

	v1 := router.Group("/api/v1")

	v1.POST("users", server.register)
	v1.POST("users/login", server.login)
	v1.POST("users/refresh", server.refresh)

	auth := v1.Group("/").Use(authMiddleware(tokenMaker))

	auth.GET("health", server.health)
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

func (s *Server) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// corsMiddleware 允許指定來源的跨域請求（開發時前端 origin）。
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	originSet := make(map[string]bool)
	for _, o := range allowedOrigins {
		originSet[o] = true
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if originSet[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
