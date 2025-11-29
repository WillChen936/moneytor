package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	TransactionTypeExpense  int16 = 1
	TransactionTypeIncome   int16 = 2
	TransactionTypeTransfer int16 = 3
)

func (server *Server) listTransactionTypes(ctx *gin.Context) {
	transcationTypes, err := server.store.ListTransactionTypes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transcationTypes)
}
