package api

import (
	db "moneytor/database/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	queries db.Querier
	router  *gin.Engine
}

func NewServer(queries db.Querier) *Server {
	server := &Server{
		queries: queries,
	}

	router := gin.Default()
	v1Routes := router.Group("/api/v1")

	v1Routes.GET("health", server.health)
	v1Routes.POST("accounts", server.createAccount)
	v1Routes.POST("categories", server.createCategory)

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
