package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) listCurrencies(ctx *gin.Context) {
	currencies, err := server.store.ListCurrencies(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, currencies)
}
