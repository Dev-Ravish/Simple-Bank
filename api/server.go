package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server starts HTTP request for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer starts a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/create-account", server.createAccount)
	router.GET("/account/:id", server.getAccount)

	server.router = router
	return server
}

// Start runs the HTTP server on the given input address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err}
}
