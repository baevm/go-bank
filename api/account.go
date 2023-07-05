package api

import (
	"database/sql"
	"errors"
	db "go-bank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	acc, err := s.db.CreateAccount(ctx, args)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(pqErr))
				return
			default:
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) GetAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	acc, err := s.db.GetAccount(ctx, req.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, "not found")
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, acc)
}

type listAccountsRequest struct {
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=5,max=20"`
}

func (s *Server) ListAccounts(ctx *gin.Context) {
	var req listAccountsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accounts, err := s.db.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
