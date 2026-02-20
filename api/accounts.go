package api

import (
	db "moneytor/database/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Name       string `json:"name" binding:"required"`
	CurrencyID int16  `json:"currencyID" binding:"required,gt=0"`
	Balance    int64  `json:"balance" binding:"required,gt=0"`
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
		Balance:    req.Balance,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.ForeignKeyViolation {
			ctx.JSON(http.StatusUnprocessableEntity, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
