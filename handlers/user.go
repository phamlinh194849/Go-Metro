package handlers

import (
	"go-metro/config"
	"go-metro/models"
	"go-metro/utils"

	"github.com/gin-gonic/gin"
)

// Register handles POST /auth/register
// @Summary Register a new user
// @Description Register a new user account with username, password, email, and full name
// @Tags auth
// @Accept json
// @Produce json
//
//	@Param request body object true "User registration data" schema={
//	  "type": "object",
//	  "required": ["username", "password", "email", "full_name"],
//	  "properties": {
//	    "username": {"type": "string", "minLength": 1},
//	    "password": {"type": "string", "minLength": 6},
//	    "email": {"type": "string", "format": "email"},
//	    "full_name": {"type": "string", "minLength": 1}
//	  }
//	}
//
// @Success 201 {object} utils.Response{data=object{user=models.User,token=string}} "User registered successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error or user already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
		FullName string `json:"full_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "username already exists")
		return
	}

	// Check if email already exists
	if err := config.DB.Where("email = ?", request.Email).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "email already exists")
		return
	}

	// Create new user
	user := models.User{
		Username: request.Username,
		Password: utils.HashPassword(request.Password),
		Email:    request.Email,
		FullName: request.FullName,
		Role:     "user", // Default role
		Status:   "active",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "failed to create user")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.InternalServerError(c, "failed to generate token")
		return
	}

	utils.SuccessResponse(c, 201, "user registered successfully", gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles POST /auth/login
// @Summary User login
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
//
//	@Param request body object true "Login credentials" schema={
//	  "type": "object",
//	  "required": ["username", "password"],
//	  "properties": {
//	    "username": {"type": "string"},
//	    "password": {"type": "string"}
//	  }
//	}
//
// @Success 200 {object} utils.Response{data=object{user=models.User,token=string}} "Login successful"
// @Failure 400 {object} utils.Response "Bad request - invalid credentials or inactive account"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Find user by username
	var user models.User
	if err := config.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		utils.BadRequest(c, "invalid username or password")
		return
	}

	// Check password
	if user.Password != utils.HashPassword(request.Password) {
		utils.BadRequest(c, "invalid username or password")
		return
	}

	// Check if user is active
	if user.Status != "active" {
		utils.BadRequest(c, "account is inactive")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.InternalServerError(c, "failed to generate token")
		return
	}

	utils.SuccessResponse(c, 200, "login successful", gin.H{
		"user":  user,
		"token": token,
	})
}

// GetProfile handles GET /user/profile
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.User} "User profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Router /user/profile [get]
func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	utils.SuccessResponse(c, 200, "", user)
}

// UpdateProfile handles PUT /user/profile
// @Summary Update user profile
// @Description Update current user's email and full name
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
//
//	@Param request body object true "Profile update data" schema={
//	  "type": "object",
//	  "properties": {
//	    "email": {"type": "string", "format": "email"},
//	    "full_name": {"type": "string"}
//	  }
//	}
//
// @Success 200 {object} utils.Response{data=models.User} "Profile updated successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error or email already exists"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /user/profile [put]
func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	// Check if user exists
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	// Bind update data
	var request struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Update user
	if request.Email != "" {
		// Check if email already exists for other user
		var existingUser models.User
		if err := config.DB.Where("email = ? AND id != ?", request.Email, userID).First(&existingUser).Error; err == nil {
			utils.BadRequest(c, "email already exists")
			return
		}
		user.Email = request.Email
	}

	if request.FullName != "" {
		user.FullName = request.FullName
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "failed to update profile")
		return
	}

	utils.SuccessResponse(c, 200, "profile updated successfully", user)
}

// ChangePassword handles PUT /user/password
// @Summary Change user password
// @Description Change current user's password
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
//
//	@Param request body object true "Password change data" schema={
//	  "type": "object",
//	  "required": ["old_password", "new_password"],
//	  "properties": {
//	    "old_password": {"type": "string"},
//	    "new_password": {"type": "string", "minLength": 6}
//	  }
//	}
//
// @Success 200 {object} utils.Response "Password changed successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error or incorrect old password"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /user/password [put]
func ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	// Check if user exists
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	// Bind request data
	var request struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Verify old password
	if user.Password != utils.HashPassword(request.OldPassword) {
		utils.BadRequest(c, "incorrect old password")
		return
	}

	// Update password
	user.Password = utils.HashPassword(request.NewPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "failed to change password")
		return
	}

	utils.SuccessResponse(c, 200, "password changed successfully", nil)
}

// Admin: Get all users
// @Summary Get all users (Admin only)
// @Description Retrieve all users in the system (Admin access required)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]models.User} "Users retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch users")
		return
	}

	utils.SuccessResponse(c, 200, "", users)
}

// Admin: Get user by ID
// @Summary Get user by ID (Admin only)
// @Description Retrieve a specific user by their ID (Admin access required)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response{data=models.User} "User retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "User not found"
// @Router /admin/users/{id} [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	utils.SuccessResponse(c, 200, "", user)
}

// Admin: Update user
// @Summary Update user (Admin only)
// @Description Update a specific user's information (Admin access required)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
//
//	@Param request body object true "User update data" schema={
//	  "type": "object",
//	  "properties": {
//	    "email": {"type": "string", "format": "email"},
//	    "full_name": {"type": "string"},
//	    "role": {"type": "string", "enum": ["user", "admin"]},
//	    "status": {"type": "string", "enum": ["active", "inactive"]}
//	  }
//	}
//
// @Success 200 {object} utils.Response{data=models.User} "User updated successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error or email already exists"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Check if user exists
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	// Bind update data
	var request struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Update user
	if request.Email != "" {
		// Check if email already exists for other user
		var existingUser models.User
		if err := config.DB.Where("email = ? AND id != ?", request.Email, id).First(&existingUser).Error; err == nil {
			utils.BadRequest(c, "email already exists")
			return
		}
		user.Email = request.Email
	}

	if request.FullName != "" {
		user.FullName = request.FullName
	}

	if request.Role != "" {
		if request.Role != "admin" && request.Role != "user" {
			utils.BadRequest(c, "invalid role")
			return
		}
		user.Role = request.Role
	}

	if request.Status != "" {
		if request.Status != "active" && request.Status != "inactive" {
			utils.BadRequest(c, "invalid status")
			return
		}
		user.Status = request.Status
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "failed to update user")
		return
	}

	utils.SuccessResponse(c, 200, "user updated successfully", user)
}

// Admin: Delete user
// @Summary Delete user (Admin only)
// @Description Delete a specific user from the system (Admin access required)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response "User deleted successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.InternalServerError(c, "failed to delete user")
		return
	}

	utils.SuccessResponse(c, 200, "user deleted successfully", nil)
}
