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

	// 允許前端開發時的 origin（Live Server 常見埠）
	allowedOrigins := []string{
		"http://localhost:5500",
		"http://127.0.0.1:5500",
	}
	router.Use(corsMiddleware(allowedOrigins))

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
