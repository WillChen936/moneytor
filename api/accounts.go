package api

import (
	db "moneytor/database/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Name           string  `json:"Name"`
	CurrencyID     int16   `json:"CurrencyID"`
	InitialBalance Decimal `json:"InitialBalance"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	// initialBalance, err := decimal.NewFromString(req.InitialBalance)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errResponse(err))
	// 	return
	// }

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
