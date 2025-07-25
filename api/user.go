package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	db "whaleWake/db/sqlc"
	"whaleWake/util"
)

// createUserRequest defines the payload for creating a new user.
// Fields:
// - UserName: required username.
// - Email: required, must be a valid email.
// - Password: required, 8-32 characters.
type createUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type userResponse struct {
	UserName   string `json:"user_name"`
	Email      string `json:"email"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	VerifiedAt string `json:"verified_at"`
}

// CreateUser handles POST /users to create a new user.
// Validates input, checks for duplicates, and inserts into the database.
// Returns 400 for bad input, 500 for server errors, 200 for success.
func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	passHashed, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		UserName: req.UserName,
		Email:    req.Email,
		Password: passHashed,
	}

	_, err = server.store.GetUserByEmail(ctx, arg.Email)
	if err == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("User already exists")))
		return
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := userResponse{
		UserName:  user.UserName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// GetUser handles GET /users/:id to retrieve a user by UUID.
// Validates UUID and fetches user from the database.
// Returns 400 for bad UUID, 404 if not found, 500 for server errors, 200 for success.
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

	userResponse := userResponse{
		UserName:   user.UserName,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02 15:04:05"),
		VerifiedAt: user.VerifiedAt.Time.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// listUsersRequest defines query parameters for paginated user listing.
// Fields:
// - PageID: required, >= 1.
// - PageSize: required, >= 1.
type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

// ListUser handles GET /users to list users with pagination.
// Validates query params and fetches users from the database.
// Returns 400 for bad params, 500 for server errors, 200 for success.
func (server *Server) ListUser(ctx *gin.Context) {
	var req listUsersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)

	var usersResponse []userResponse
	for _, user := range users {
		usersResponse = append(usersResponse, userResponse{
			UserName:   user.UserName,
			Email:      user.Email,
			CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  user.UpdatedAt.Format("2006-01-02 15:04:05"),
			VerifiedAt: user.VerifiedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, usersResponse)
}

// DeleteUser handles DELETE /users/:id to delete a user by UUID.
// Validates UUID and deletes user from the database.
// Returns 400 for bad UUID, 500 for server errors, 200 for success.
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

	userResponse := userResponse{
		UserName:   user.UserName,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02 15:04:05"),
		VerifiedAt: user.VerifiedAt.Time.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// updateUserRequest defines the payload for updating a user.
// Fields:
// - ID: required UUID of the user.
// - UserName, Email, Password: optional new values.
type updateUserRequest struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	UserName string    `json:"user_name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

// UpdateUser handles PUT /users to update user details.
// Validates input and updates user in the database.
// Returns 400 for bad input, 500 for server errors, 200 for success.
func (server *Server) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	passHashed, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:       req.ID,
		UserName: req.UserName,
		Email:    req.Email,
		Password: passHashed,
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := userResponse{
		UserName:   user.UserName,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02 15:04:05"),
		VerifiedAt: user.VerifiedAt.Time.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// createUserTxRequest defines the payload for transactional user creation.
// Includes user, profile, and role fields.
// All fields are required except for role, which is set internally.
type createUserTxRequest struct {
	UserName      string `json:"user_name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=8,max=32"`
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	BusinessName  string `json:"business_name" binding:"required"`
	StreetAddress string `json:"street_address" binding:"required"`
	City          string `json:"city" binding:"required"`
	State         string `json:"state" binding:"required"`
	Zip           string `json:"zip" binding:"required"`
	CountryCode   string `json:"country_code" binding:"required"`
}

type createUserTxResponse struct {
	UserName      string `json:"user_name"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BusinessName  string `json:"business_name"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zip           string `json:"zip"`
	CountryCode   string `json:"country_code"`
	RoleID        int32  `json:"role_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	VerifiedAt    string `json:"verified_at"`
}

// CreateUserTx handles POST /users/tx for transactional user creation.
// Creates user, profile, and role in a single transaction.
// Returns 400 for bad input, 500 for server errors, 200 for success.
func (server *Server) CreateUserTx(ctx *gin.Context) {
	var req createUserTxRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	passHashed, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userParams := db.CreateUserParams{
		UserName: req.UserName,
		Email:    req.Email,
		Password: passHashed,
	}

	profileParams := db.CreateUserProfileParams{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		BusinessName:  req.BusinessName,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Zip:           req.Zip,
		CountryCode:   req.CountryCode,
	}

	roleParams := db.CreateUserRoleParams{
		RoleID: 1,
	}

	_, err = server.store.GetUserByEmail(ctx, userParams.Email)
	if err == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("User already exists")))
		return
	}

	userWithProfileAndRole, err := server.store.CreateUserWithProfileAndRoleTx(ctx, userParams, profileParams, roleParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := createUserTxResponse{
		UserName:      userWithProfileAndRole.User.UserName,
		Email:         userWithProfileAndRole.User.Email,
		FirstName:     userWithProfileAndRole.UserProfile.FirstName,
		LastName:      userWithProfileAndRole.UserProfile.LastName,
		BusinessName:  userWithProfileAndRole.UserProfile.BusinessName,
		StreetAddress: userWithProfileAndRole.UserProfile.StreetAddress,
		City:          userWithProfileAndRole.UserProfile.City,
		State:         userWithProfileAndRole.UserProfile.State,
		Zip:           userWithProfileAndRole.UserProfile.Zip,
		CountryCode:   userWithProfileAndRole.UserProfile.CountryCode,
		RoleID:        userWithProfileAndRole.UserRole.RoleID,
		CreatedAt:     userWithProfileAndRole.User.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     userWithProfileAndRole.User.UpdatedAt.Format("2006-01-02 15:04:05"),
		VerifiedAt:    userWithProfileAndRole.User.VerifiedAt.Time.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// getUserTxRequest defines the payload for transactional user retrieval.
// Field:
// - ID: required UUID of the user.
type getUserTxRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

// GetUserTx handles GET /users/tx/:id for transactional user retrieval.
// Fetches user with profile and role in a single transaction.
// Returns 400 for bad UUID, 500 for server errors, 200 for success.
func (server *Server) GetUserTx(ctx *gin.Context) {
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

	userWithProfileAndRole, err := server.store.GetUserWithProfileAndRoleTX(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := createUserTxResponse{
		UserName:      userWithProfileAndRole.User.UserName,
		Email:         userWithProfileAndRole.User.Email,
		FirstName:     userWithProfileAndRole.UserProfile.FirstName,
		LastName:      userWithProfileAndRole.UserProfile.LastName,
		BusinessName:  userWithProfileAndRole.UserProfile.BusinessName,
		StreetAddress: userWithProfileAndRole.UserProfile.StreetAddress,
		City:          userWithProfileAndRole.UserProfile.City,
		State:         userWithProfileAndRole.UserProfile.State,
		Zip:           userWithProfileAndRole.UserProfile.Zip,
		CountryCode:   userWithProfileAndRole.UserProfile.CountryCode,
		RoleID:        userWithProfileAndRole.UserRole.RoleID,
		CreatedAt:     userWithProfileAndRole.User.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     userWithProfileAndRole.User.UpdatedAt.Format("2006-01-02 15:04:05"),
		VerifiedAt:    userWithProfileAndRole.User.VerifiedAt.Time.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// DeleteUserTx handles DELETE /users/tx/:id for transactional user deletion.
// Validates UUID, checks store initialization, and deletes user with profile and role in a single transaction.
// Returns 400 for bad UUID, 500 for server errors, 200 for success.
func (server *Server) DeleteUserTx(ctx *gin.Context) {
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

	userWithProfileAndRole, err := server.store.DeleteUserWithProfileAndRoleTX(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := createUserTxResponse{
		UserName:      userWithProfileAndRole.User.UserName,
		Email:         userWithProfileAndRole.User.Email,
		FirstName:     userWithProfileAndRole.UserProfile.FirstName,
		LastName:      userWithProfileAndRole.UserProfile.LastName,
		BusinessName:  userWithProfileAndRole.UserProfile.BusinessName,
		StreetAddress: userWithProfileAndRole.UserProfile.StreetAddress,
		City:          userWithProfileAndRole.UserProfile.City,
		State:         userWithProfileAndRole.UserProfile.State,
		Zip:           userWithProfileAndRole.UserProfile.Zip,
		CountryCode:   userWithProfileAndRole.UserProfile.CountryCode,
		RoleID:        userWithProfileAndRole.UserRole.RoleID,
		CreatedAt:     userWithProfileAndRole.User.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     userWithProfileAndRole.User.UpdatedAt.Format("2006-01-02 15:04:05"),
		VerifiedAt:    userWithProfileAndRole.User.VerifiedAt.Time.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// createUserTxRequest defines the payload for transactional user creation.
// Includes user, profile, and role fields.
// All fields are required except for role, which is set internally.
type updateUserTxRequest struct {
	ID            uuid.UUID `json:"id" binding:"required"`
	UserName      string    `json:"user_name"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	BusinessName  string    `json:"business_name"`
	StreetAddress string    `json:"street_address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Zip           string    `json:"zip"`
	CountryCode   string    `json:"country_code"`
	RoleID        int32     `json:"role_id"`
}

func (server *Server) UpdateUserTx(ctx *gin.Context) {
	var req updateUserTxRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if server.store == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("store not initialized")))
		return
	}

	passHashed, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updateUserParams := db.UpdateUserParams{
		ID:       req.ID,
		UserName: req.UserName,
		Email:    req.Email,
		Password: passHashed,
	}

	updateProfileParams := db.UpdateUserProfileParams{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		BusinessName:  req.BusinessName,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Zip:           req.Zip,
		CountryCode:   req.CountryCode,
	}

	updateRoleParams := db.UpdateUserRoleParams{
		RoleID: req.RoleID,
	}

	updatedUserWithProfileAndRole, err := server.store.UpdateUserWithProfileAndRoleTX(ctx, updateUserParams, updateProfileParams, updateRoleParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := createUserTxResponse{
		UserName:      updatedUserWithProfileAndRole.User.UserName,
		Email:         updatedUserWithProfileAndRole.User.Email,
		FirstName:     updatedUserWithProfileAndRole.UserProfile.FirstName,
		LastName:      updatedUserWithProfileAndRole.UserProfile.LastName,
		BusinessName:  updatedUserWithProfileAndRole.UserProfile.BusinessName,
		StreetAddress: updatedUserWithProfileAndRole.UserProfile.StreetAddress,
		City:          updatedUserWithProfileAndRole.UserProfile.City,
		State:         updatedUserWithProfileAndRole.UserProfile.State,
		Zip:           updatedUserWithProfileAndRole.UserProfile.Zip,
		CountryCode:   updatedUserWithProfileAndRole.UserProfile.CountryCode,
		RoleID:        updatedUserWithProfileAndRole.UserRole.RoleID,
		CreatedAt:     updatedUserWithProfileAndRole.User.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     updatedUserWithProfileAndRole.User.UpdatedAt.Format("2006-01-02 15:04:05"),
		VerifiedAt:    updatedUserWithProfileAndRole.User.VerifiedAt.Time.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, userResponse)
}
