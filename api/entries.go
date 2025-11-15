package api

import (
	db "moneytor/database/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createEntryRequest struct {
	Name       string  `json:"name" binding:"required,min=1,max=50"`
	Note       string  `json:"note" binding:"required,min=1,max=200"`
	AccountID  int64   `json:"accountId" binding:"required"`
	CategoryID int64   `json:"categoryId" binding:"required"`
	Amount     Decimal `json:"amount" binding:"required"`
}

func (server *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateEntryParams{
		Name:       req.Name,
		Note:       req.Note,
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount.Decimal,
	}

	entry, err := server.queries.CreateEntry(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.ForeignKeyViolation {
			ctx.JSON(http.StatusUnprocessableEntity, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
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

		entries, err = server.queries.ListEntries(ctx, arg)
	} else {
		arg := db.ListEntriesByAccountIDParams{
			AccountID: req.AccountID,
			Limit:     req.PageSize,
			Offset:    offset,
		}

		entries, err = server.queries.ListEntriesByAccountID(ctx, arg)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}
