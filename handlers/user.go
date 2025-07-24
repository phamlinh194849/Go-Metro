package handlers

import (
	"go-metro/config"
	"go-metro/consts"
	"go-metro/models"
	"go-metro/utils"

	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

type LoginReq struct {
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRes struct {
	Email    string `json:"email"`
	Role     string `json:"role"`
	Username string `json:"username"`
}

// Bind update data
type UpdateInfoReq struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

// Bind request data
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type AdminUpdateInfoReq struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

// Ok
// Register handles POST /auth/register
// @Summary Register a new user
// @Description Register a new user account with username, password, email, and full name
// @Tags auth
// @Accept json
// @Produce json
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var request RegisterReq

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "Tên người dùng đã tồn")
		return
	}

	// Check if email already exists
	if err := config.DB.Where("email = ?", request.Email).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "Email đã tồn tại")
		return
	}

	// Create new user
	user := models.User{
		Username: request.Username,
		Password: utils.HashPassword(request.Password),
		Email:    request.Email,
		FullName: request.FullName,
		Role:     2, // Default role user
		Status:   "active",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "Lỗi khi tạo tài khoản")
		return
	}

	utils.SuccessResponse(c, 201, "Đã tạo tài khoản thành công", gin.H{
		"user": user,
	})
}

// Ok
// Login handles POST /auth/login
// @Summary User login
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "Login credentials" schema={
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var request LoginReq

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Find user by username
	var user models.User
	if err := config.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		utils.BadRequest(c, "Email không đúng")
		return
	}

	// Check password
	if user.Password != utils.HashPassword(request.Password) {
		utils.BadRequest(c, "Sai mật khẩu")
		return
	}

	// Check if user is active
	if user.Status != "active" {
		utils.BadRequest(c, "Tài khoản bị khóa")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, int(user.Role))
	if err != nil {
		utils.InternalServerError(c, "Lỗi trong quá trình tạo token ")
		return
	}

	login := LoginRes{
		Email:    user.Email,
		Role:     user.Role.ToText(),
		Username: user.Username,
	}

	utils.SuccessResponse(c, 200, "login successful", gin.H{
		"user":  login,
		"token": token,
	})
}

// Ok
// GetProfile handles GET /user/profile
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /user/profile [get]
func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "Thông tin không tồn tại")
		return
	}

	utils.SuccessResponse(c, 200, "", user)
}

// Ok
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
// @Failure 404 {object} utils.Response "Thông tin không tồn tại"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /user/profile [put]
func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	// Check if user exists
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "Thông tin không tồn tại")
		return
	}

	// Bind update data
	request := UpdateInfoReq{}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Update user
	if request.Username != "" {
		var existingUser models.User
		if err := config.DB.Where("username = ? AND id != ?", request.Username, userID).First(&existingUser).Error; err == nil {
			utils.BadRequest(c, "username already exists")
			return
		}
		user.Username = request.Username
	}

	if request.FullName != "" {
		user.FullName = request.FullName
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "Lỗi khi cập nhật thông tin người dùng")
		return
	}

	utils.SuccessResponse(c, 200, "Cập nhật thành công", user)
}

// Ok
// ChangePassword handles PUT /user/password
// @Summary Change user password
// @Description Change current user's password
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Password change data" schema={
// @Router /user/password [put]
func ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	// Check if user exists
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "Thông tin không tồn tại")
		return
	}

	// Bind request data
	request := ChangePasswordReq{}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Verify old password
	if user.Password != utils.HashPassword(request.OldPassword) {
		utils.BadRequest(c, "Mật khẩu cũ không đúng")
		return
	}

	// Update password
	user.Password = utils.HashPassword(request.NewPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "failed to change password")
		return
	}

	utils.SuccessResponse(c, 200, "Đổi mật khẩu thành công", nil)
}

// Ok
// Admin: Get all users
// @Summary Get all users (Admin only)
// @Description Retrieve all users in the system (Admin access required)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /admin/users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		utils.InternalServerError(c, "Lỗi khi lấy danh sách người dùng")
		return
	}

	utils.SuccessResponse(c, 200, "", users)
}

// Ok
// Admin: Get user by ID
// @Summary Get user by ID (Admin only)
// @Description Retrieve a specific user by their ID (Admin access required)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Router /admin/users/{id} [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "Thông tin không tồn tại")
		return
	}

	utils.SuccessResponse(c, 200, "", user)
}

// Ok
// Admin: Update user
// @Summary Update user (Admin only)
// @Description Update a specific user's information (Admin access required)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body object true "User update data" schema={
// @Router /admin/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Check if user exists
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "Thông tin không tồn tại")
		return
	}

	// Bind update data
	request := AdminUpdateInfoReq{}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Update user
	if request.Email != "" {
		var existingUser models.User
		if err := config.DB.Where("email = ? AND id != ?", request.Email, id).First(&existingUser).Error; err == nil {
			utils.BadRequest(c, "Email đã tồn tại")
			return
		}
		user.Email = request.Email
	}
	if request.Username != "" {
		var existingUser models.User
		if err := config.DB.Where("username = ? AND id != ?", request.Username, id).First(&existingUser).Error; err == nil {
			utils.BadRequest(c, "Tên người dùng đã tồn tại")
			return
		}
		user.Username = request.Username
	}
	if request.FullName != "" {
		user.FullName = request.FullName
	}

	if request.Role != "" {
		roleMap := map[string]consts.Role{
			"admin": consts.AdminRole,
			"user":  consts.UserRole,
			"staff": consts.StaffRole,
		}

		roleValue, exists := roleMap[request.Role]
		if !exists {
			utils.BadRequest(c, "Vai trò không đúng")
			return
		}

		user.Role = roleValue
	}

	if request.Status != "" {
		if request.Status != "active" && request.Status != "inactive" {
			utils.BadRequest(c, "Không xác định")
			return
		}
		user.Status = request.Status
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "failed to update user")
		return
	}

	utils.SuccessResponse(c, 200, "Cập nhật thông tin người dùng thành công", user)
}

// Ok
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
// @Failure 404 {object} utils.Response "Thông tin không tồn tại"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "Thông tin không tồn tại")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.InternalServerError(c, "failed to delete user")
		return
	}

	utils.SuccessResponse(c, 200, "Xóa thành công", nil)
}
