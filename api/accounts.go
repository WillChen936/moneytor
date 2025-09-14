package api

import (
	db "moneytor/database/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type createAccountRequest struct {
	Name           string          `json:"Name"`
	CurrencyID     int16           `json:"CurrencyID"`
	InitailBalance decimal.Decimal `json:"InitailBalance"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	log.Info().Msgf("req.Balance = %s", req.InitailBalance.String())

	arg := db.CreateAccountParams{
		Name:       req.Name,
		CurrencyID: req.CurrencyID,
		Balance:    req.InitailBalance,
	}

	log.Info().Msgf("arg.Balance = %s", arg.Balance.String())

	account, err := server.queries.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	log.Info().Msgf("account.Balance = %s", account.Balance.String())

	ctx.JSON(http.StatusOK, account)
}
