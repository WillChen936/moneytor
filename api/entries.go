package api

import (
	"fmt"
	db "moneytor/database/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createEntryRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=50"`
	Note       string `json:"note"`
	AccountID  int64  `json:"accountId" binding:"required"`
	CategoryID int64  `json:"categoryId" binding:"required"`
	Amount     int64  `json:"amount" binding:"required"`
}

func (server *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	category, err := server.store.GetCategory(ctx, req.CategoryID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}

	amount, err := ResolverEntryAmount(category.TransactionTypeID, req.Amount)
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

type getEntriesRequest struct {
	AccountID int64 `form:"account_id"`
	PageID    int32 `form:"page_id,default=1" binding:"min=1"`
	PageSize  int32 `form:"page_size,default=5" binding:"min=5,max=10"`
}

func (server *Server) getEntries(ctx *gin.Context) {
	var req getEntriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var entries []db.Entry
	var err error
	offset := (req.PageID - 1) * req.PageSize
	if req.AccountID == 0 {
		arg := db.ListEntriesParams{
			Limit:  req.PageSize,
			Offset: offset,
		}

		entries, err = server.store.ListEntries(ctx, arg)
	} else {
		arg := db.ListEntriesByAccountIDParams{
			AccountID: req.AccountID,
			Limit:     req.PageSize,
			Offset:    offset,
		}

		entries, err = server.store.ListEntriesByAccountID(ctx, arg)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

func ResolverEntryAmount(transactionTypeID int16, rawAmount int64) (int64, error) {
	switch transactionTypeID {
	case TransactionTypeExpense:
		return -1 * rawAmount, nil
	case TransactionTypeIncome:
		return 1 * rawAmount, nil
	default:
		return 0, fmt.Errorf("invalid transaction type: %d", transactionTypeID)
	}
}
