package api

import (
	db "moneytor/database/sqlc"
	"moneytor/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=50"`
	CurrencyID int16  `json:"currencyId" binding:"required,gt=0"`
	Balance    int64  `json:"balance" binding:"required,gte=0,lte=999999999999"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		UserID:     payload.UserID,
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

type listAccountsRequest struct {
	PageID   int32 `form:"pageId,default=1" binding:"min=1"`
	PageSize int32 `form:"pageSize,default=5" binding:"min=1,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		UserID: payload.UserID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
