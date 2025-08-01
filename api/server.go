package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "whaleWake/db/sqlc"
	"whaleWake/token"
	"whaleWake/util"
)

// Server serves HTTP requests for the application.
type Server struct {
	config     util.Config // Configuration settings for the server.
	store      db.Store    // Database store for executing queries.
	tokenMaker token.Maker // Token maker for generating and validating tokens.
	router     *gin.Engine // HTTP router for handling API routes.
}

// NewServer creates a new Server instance and sets up the routes.
// Parameters:
// - store: A pointer to the db.Store instance for database operations.
// Returns:
// - A pointer to the newly created Server instance.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	// Assign the configured router to the server.
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// Basic User CRUD Routes
	router.POST("/users", server.CreateUser)       // Create a new user.
	router.GET("/users/:id", server.GetUser)       // Retrieve a user by ID.
	router.GET("/users", server.ListUser)          // List all users.
	router.DELETE("/users/:id", server.DeleteUser) // Delete a user by ID.
	router.PUT("/users", server.UpdateUser)        // Update user details.
	router.POST("/users/login", server.LoginUser)  // User login route.

	// User Transaction (TX) Routes
	router.POST("/usertx", server.CreateUserTx)       // Create a user transaction.
	router.GET("/usertx/:id", server.GetUserTx)       // Retrieve user transactions.
	router.DELETE("/usertx/:id", server.DeleteUserTx) // Delete user transactions.
	router.PUT("/usertx", server.UpdateUserTx)        // Update user transactions.
	server.router = router
}

// Start runs the HTTP server on the specified address.
// Parameters:
// - address: The address (host:port) to bind the server to.
// Returns:
// - An error if the server fails to start.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse formats an error message as a JSON response.
// Parameters:
// - err: The error to be formatted.
// Returns:
// - A gin.H map containing the error message.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
