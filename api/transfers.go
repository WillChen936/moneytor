package api

import (
	"fmt"
	db "moneytor/database/sqlc"
	"moneytor/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	Name          string `json:"name" binding:"required,min=1,max=50"`
	Note          string `json:"note"`
	FromAccountID int64  `json:"fromAccountId" binding:"required,gt=0"`
	ToAccountID   int64  `json:"toAccountId" binding:"required,gt=0"`
	CategoryID    int64  `json:"categoryId" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	if req.FromAccountID == req.ToAccountID {
		ctx.JSON(http.StatusBadRequest, errResponse(fmt.Errorf("fromAccountId and toAccountId must be different")))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	category, err := server.store.GetCategory(ctx, db.GetCategoryParams{
		ID:     req.CategoryID,
		UserID: payload.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}

	if category.TransactionTypeID != TransactionTypeTransfer {
		ctx.JSON(http.StatusBadRequest, errResponse(fmt.Errorf("category is not a transfer type")))
		return
	}

	_, err = server.store.GetAccount(ctx, db.GetAccountParams{
		ID:     req.ToAccountID,
		UserID: payload.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}

	arg := db.CreateTransferTxParams{
		Name:          req.Name,
		Note:          req.Note,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		CategoryID:    req.CategoryID,
		Amount:        req.Amount,
	}

	result, err := server.store.CreateTransferTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.ForeignKeyViolation {
			ctx.JSON(http.StatusUnprocessableEntity, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}
