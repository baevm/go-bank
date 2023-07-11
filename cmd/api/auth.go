package api

import (
	"database/sql"
	"errors"
	db "go-bank/db/sqlc"
	"go-bank/internal/password"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID          uuid.UUID `json:"sessionId"`
	Username           string    `json:"username"`
	AccessToken        string    `json:"access_token"`
	AccessTokenExpire  time.Time `json:"access_token_expire"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpire time.Time `json:"refresh_token_expire"`
}

func (s *HTTPServer) Login(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.db.GetUser(ctx, req.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	err = password.Check(user.HashedPass, req.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.Create(req.Username, s.cfg.ACCESS_TOKEN_DURATION)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.Create(req.Username, s.cfg.REFRESH_TOKEN_DURATION)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := s.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, &loginUserResponse{
		SessionID:          session.ID,
		Username:           req.Username,
		AccessToken:        accessToken,
		AccessTokenExpire:  accessPayload.ExpiresAt.Time,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: refreshPayload.ExpiresAt.Time,
	})
}

type refreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type refreshAccessTokenResponse struct {
	AccessToken       string    `json:"access_token"`
	AccessTokenExpire time.Time `json:"access_token_expire"`
}

func (s *HTTPServer) RefreshAccessToken(ctx *gin.Context) {
	var req refreshAccessTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := s.tokenMaker.Verify(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := s.db.GetSession(ctx, payload.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	if session.IsBlocked {
		err := errors.New("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != payload.Username {
		err := errors.New("bad session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := errors.New("bad session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := errors.New("bad session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.Create(payload.Username, s.cfg.ACCESS_TOKEN_DURATION)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, &refreshAccessTokenResponse{
		AccessToken:       accessToken,
		AccessTokenExpire: accessPayload.ExpiresAt.Time,
	})
}
