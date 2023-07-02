package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "go-bank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) CreateTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	if !s.isValidAccountCurrency(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !s.isValidAccountCurrency(ctx, req.ToAccountID, req.Currency) {
		return
	}

	transfer, err := s.db.TransferTx(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (s *Server) isValidAccountCurrency(ctx *gin.Context, accountId int64, currency string) bool {
	acc, err := s.db.GetAccount(ctx, accountId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, "not found")
			return false
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return false
		}
	}

	if acc.Currency != currency {
		err = fmt.Errorf("account currency mismatch")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
