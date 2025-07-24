package routes

import (
	"go-metro/handlers"
	"go-metro/utils"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine) {
	// History routes
	historyGroup := r.Group("/history")
	{
		historyGroup.POST("", handlers.CreateHistory)
		historyGroup.GET("", handlers.GetHistories)
		historyGroup.GET("/:id", handlers.GetHistoryByID)
	}

	// Card routes
	cardGroup := r.Group("/card")
	{
		cardGroup.POST("", handlers.CreateCard)                     // Tạo card mới
		cardGroup.GET("", handlers.GetCards)                        // Lấy danh sách tất cả cards
		cardGroup.GET("/:id", handlers.GetCardByID)                 // Lấy card theo ID
		cardGroup.GET("/cardid/:rf_id", handlers.GetCardByCardID)   // Lấy card theo CardID
		cardGroup.PUT("/:id", handlers.UpdateCard)                  // Cập nhật card
		cardGroup.DELETE("/:id", handlers.DeleteCard)               // Xóa card
		cardGroup.POST("/:id/topup", handlers.TopUpCard)            // Nạp tiền vào card
		cardGroup.GET("/user/:user_id", handlers.GetCardsByUser)    // Lấy cards theo user_id
		cardGroup.GET("/status/:status", handlers.GetCardsByStatus) // Lấy cards theo status
	}

	// Auth routes (public)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", handlers.Register) // Đăng ký
		authGroup.POST("/login", handlers.Login)       // Đăng nhập
	}

	// User routes (require authentication)
	userGroup := r.Group("/user")
	userGroup.Use(utils.AuthMiddleware())
	{
		userGroup.GET("/profile", handlers.GetProfile)      // Xem profile
		userGroup.PUT("/profile", handlers.UpdateProfile)   // Cập nhật profile
		userGroup.PUT("/password", handlers.ChangePassword) // Đổi mật khẩu
	}

	// Admin routes (require admin role)
	adminGroup := r.Group("/admin")
	adminGroup.Use(utils.AuthMiddleware(), utils.AdminMiddleware())
	{
		adminGroup.GET("/users", handlers.GetAllUsers)       // Lấy tất cả users
		adminGroup.GET("/users/:id", handlers.GetUserByID)   // Lấy user theo ID
		adminGroup.PUT("/users/:id", handlers.UpdateUser)    // Cập nhật user
		adminGroup.DELETE("/users/:id", handlers.DeleteUser) // Xóa user
	}

	// Health check route
	// @Summary Health check
	// @Description Check if the API is running
	// @Tags system
	// @Accept json
	// @Produce json
	// @Success 200 {object} object{status=string} "API is healthy"
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
