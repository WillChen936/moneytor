package api

import (
	"errors"
	db "moneytor/database/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCategoryRequest struct {
	Name              string `json:"name" binding:"required"`
	TransactionTypeID int16  `json:"transactionTypeId" binding:"required,gt=0"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateCategoryParams{
		Name:              req.Name,
		TransactionTypeID: req.TransactionTypeID,
	}

	category, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.ForeignKeyViolation {
			ctx.JSON(http.StatusUnprocessableEntity, errResponse(errors.New("transaction type not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type listCategoriesRequest struct {
	PageID   int32 `form:"pageId,default=1" binding:"min=1"`
	PageSize int32 `form:"pageSize,default=5" binding:"min=5,max=10"`
}

func (server *Server) listCategories(ctx *gin.Context) {
	var req listCategoriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.ListCategoriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	category, err := server.store.ListCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}
