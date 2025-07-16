package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	db "whaleWake/db/sqlc"
)

type createUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	_, err := server.store.GetUserByEmail(ctx, arg.Email)
	if err == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("User already exists")))
		return
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
