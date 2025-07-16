package handlers

import (
	"go-metro/config"
	"go-metro/models"
	"go-metro/utils"

	"github.com/gin-gonic/gin"
)

// Register handles POST /auth/register
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
func GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch users")
		return
	}

	utils.SuccessResponse(c, 200, "", users)
}

// Admin: Get user by ID
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
