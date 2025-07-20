package api

import (
	"github.com/gin-gonic/gin"
	db "whaleWake/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/users", server.CreateUser)
	router.GET("/users/:id", server.GetUser)
	router.GET("/users", server.ListUser)
	router.DELETE("/users/:id", server.DeleteUser)
	router.PUT("/users", server.UpdateUser)
	router.POST("/usertx", server.CreateUserTx)

	server.router = router // Assign router to server
	//Add Routes to Router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
