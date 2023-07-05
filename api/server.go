package api

import (
	db "go-bank/db/sqlc"
	vl "go-bank/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	db     db.Store
	router *gin.Engine
}

func NewServer(db db.Store) *Server {
	server := &Server{
		db: db,
	}

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", vl.ValidCurrency)
	}

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccounts)

	router.POST("/transfers", server.CreateTransfer)

	router.POST("/users", server.CreateUser)

	server.router = router

	return server
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
