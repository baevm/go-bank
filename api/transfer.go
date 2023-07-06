package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "go-bank/db/sqlc"
	"go-bank/internal/token"
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	fromAccount, isValid := s.isValidAccountCurrency(ctx, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}

	if authPayload.Username != fromAccount.Owner {
		err := errors.New("account not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, isValid = s.isValidAccountCurrency(ctx, req.ToAccountID, req.Currency)
	if !isValid {
		return
	}

	transfer, err := s.db.TransferTx(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (s *Server) isValidAccountCurrency(ctx *gin.Context, accountId int64, currency string) (db.Accounts, bool) {
	acc, err := s.db.GetAccount(ctx, accountId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, "not found")
			return acc, false
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return acc, false
		}
	}

	if acc.Currency != currency {
		err = fmt.Errorf("account currency mismatch")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return acc, false
	}

	return acc, true
}
