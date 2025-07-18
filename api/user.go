package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (server *Server) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	user, err := server.store.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// `form` pulls off the URI query item.
type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

func (server *Server) ListUser(ctx *gin.Context) {

	var req listUsersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	users, err := server.store.ListUsers(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	user, err := server.store.DeleteUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, user)
}

type updateUserRequest struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	UserName string    `json:"user_name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (server *Server) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:       req.ID,
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// TODO: Make a CreateUserTx Handler
// TODO: Make an UpdateUserTX Handler
// TODO: Make a DeleteUserTX Handler
// TODO: Make a GetUserTX Handler
// TODO: Make a ListUserTX Handler
