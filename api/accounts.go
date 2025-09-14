package api

import (
	db "moneytor/database/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Name           string  `json:"Name" binding:"required"`
	CurrencyID     int16   `json:"CurrencyID" binding:"gt=0"`
	InitialBalance Decimal `json:"InitialBalance" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Name:       req.Name,
		CurrencyID: req.CurrencyID,
		Balance:    req.InitialBalance.Decimal,
	}

	account, err := server.queries.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
