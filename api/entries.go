package api

import (
	"fmt"
	db "moneytor/database/sqlc"
	"moneytor/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createEntryRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=50"`
	Note       string `json:"note"`
	AccountID  int64  `json:"accountId" binding:"required"`
	CategoryID int64  `json:"categoryId" binding:"required"`
	Amount     int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := ctx.MustGet(authPayloadKey).(*token.Payload)

	category, err := server.store.GetCategory(ctx, db.GetCategoryParams{
		ID:     req.CategoryID,
		UserID: payload.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}

	amount, err := resolverEntryAmount(category.TransactionTypeID, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateEntryTxParams{
		Name:       req.Name,
		Note:       req.Note,
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		Amount:     amount,
	}

	result, err := server.store.CreateEntryTx(ctx, arg)
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

type listEntriesRequest struct {
	AccountID int64 `form:"accountId" binding:"required,gt=0"`
	PageID    int32 `form:"pageId,default=1" binding:"min=1"`
	PageSize  int32 `form:"pageSize,default=5" binding:"min=1,max=10"`
}

func (server *Server) listEntries(ctx *gin.Context) {
	var req listEntriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	offset := (req.PageID - 1) * req.PageSize
	arg := db.ListEntriesByAccountIDParams{
		AccountID: req.AccountID,
		Limit:     req.PageSize,
		Offset:    offset,
	}

	entries, err := server.store.ListEntriesByAccountID(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

func resolverEntryAmount(transactionTypeID int16, rawAmount int64) (int64, error) {
	switch transactionTypeID {
	case TransactionTypeExpense:
		return -1 * rawAmount, nil
	case TransactionTypeIncome:
		return 1 * rawAmount, nil
	default:
		return 0, fmt.Errorf("invalid transaction type: %d", transactionTypeID)
	}
}
