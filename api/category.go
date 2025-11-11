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

	category, err := server.queries.CreateCategory(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.ForeignKeyViolation {
			ctx.JSON(http.StatusNotFound, errResponse(errors.New("transaction type not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}
