package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	db "whaleWake/db/sqlc"
	"whaleWake/token"
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
	ID         uuid.UUID `json:"id"`
	UserName   string    `json:"user_name"`
	Email      string    `json:"email"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
	VerifiedAt string    `json:"verified_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:         user.ID,
		UserName:   user.UserName,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02 15:04:05"),
		VerifiedAt: user.VerifiedAt.Time.Format("2006-01-02 15:04:05"),
	}
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

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		UserName: req.UserName,
		Email:    req.Email,
		Password: hashedPassword,
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

	//We're going to give the user a Role off the rip that way we can Auth roles later.
	userRoleParams := db.CreateUserRoleParams{
		UserID: user.ID,
		RoleID: 1,
	}

	_, err = server.store.CreateUserRole(ctx, userRoleParams)

	userResponse := newUserResponse(user)

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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if user.ID != authPayload.UserID && authPayload.RoleID != 3 {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("You are not authorized to view this user")))
		return
	}

	userResponse := newUserResponse(user)

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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.RoleID != 3 {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("You are not authorized to view these users")))
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
			ID:         user.ID,
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.RoleID != 3 {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("You are not authorized to delete this user")))
		return
	}

	user, err := server.store.DeleteUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	userResponse := newUserResponse(user)

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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.ID != authPayload.UserID && authPayload.RoleID != 3 {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("You are not authorized to view this user")))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:       req.ID,
		UserName: req.UserName,
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := newUserResponse(user)

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
	ID            uuid.UUID `json:"id"`
	UserName      string    `json:"user_name"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	BusinessName  string    `json:"business_name"`
	StreetAddress string    `json:"street_address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Zip           string    `json:"zip"`
	CountryCode   string    `json:"country_code"`
	RoleID        int32     `json:"role_id"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	VerifiedAt    string    `json:"verified_at"`
}

func newUserTXResponse(userWithProfileAndRole db.UserTxResult) createUserTxResponse {
	return createUserTxResponse{
		ID:            userWithProfileAndRole.User.ID,
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

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userParams := db.CreateUserParams{
		UserName: req.UserName,
		Email:    req.Email,
		Password: hashedPassword,
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

	userResponse := newUserTXResponse(userWithProfileAndRole)

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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if userWithProfileAndRole.User.ID != authPayload.UserID && authPayload.RoleID != 3 {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("You are not authorized to view this user")))
		return
	}

	userResponse := newUserTXResponse(userWithProfileAndRole)

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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.RoleID != 3 {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("You are not authorized to delete users")))
		return
	}

	userWithProfileAndRole, err := server.store.DeleteUserWithProfileAndRoleTX(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := newUserTXResponse(userWithProfileAndRole)

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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.ID != authPayload.UserID && authPayload.RoleID != 3 {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("You are not authorized to update this user")))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updateUserParams := db.UpdateUserParams{
		ID:       req.ID,
		UserName: req.UserName,
		Email:    req.Email,
		Password: hashedPassword,
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

	userWithProfileAndRole, err := server.store.UpdateUserWithProfileAndRoleTX(ctx, updateUserParams, updateProfileParams, updateRoleParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := newUserTXResponse(userWithProfileAndRole)

	ctx.JSON(http.StatusOK, userResponse)
}

// loginUserRequest defines the payload for logging in a user.
// Fields:
// - Email: required, must be a valid email.
// - Password: required, 8-32 characters.
type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid email or password")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid email or password")))
		return
	}

	userRole, err := server.store.GetUserRole(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.ID,
		int(userRole.RoleID),
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
