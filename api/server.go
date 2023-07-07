package api

import (
	"fmt"
	"go-bank/config"
	db "go-bank/db/sqlc"
	"go-bank/internal/token"
	vl "go-bank/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	db         db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	cfg        config.Config
}

func NewServer(config config.Config, db db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TOKEN_SYMMETRIC_KEY)

	if err != nil {
		return nil, fmt.Errorf("cant create server: %s", err)
	}

	server := &Server{
		db:         db,
		tokenMaker: tokenMaker,
		cfg:        config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", vl.ValidCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/users", s.CreateUser)
	router.POST("/users/login", s.Login)
	router.POST("/users/refreshToken", s.RefreshAccessToken)

	authRoutes := router.Group("/")
	{
		authRoutes.Use(AuthMiddleware(s.tokenMaker))

		authRoutes.POST("/accounts", s.CreateAccount)
		authRoutes.GET("/accounts/:id", s.GetAccount)
		authRoutes.GET("/accounts", s.ListAccounts)

		authRoutes.POST("/transfers", s.CreateTransfer)
	}

	s.router = router

	return router
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
