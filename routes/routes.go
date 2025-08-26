package routes

import (
  "go-metro/handlers"
  "go-metro/utils"

  "github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine) {
  // History routes (read-only)
  historyGroup := r.Group("/history")
  {
    historyGroup.GET("", handlers.GetHistories)
    historyGroup.GET("/:id", handlers.GetHistoryByID)
  }

  // Sell History routes (read-only)
  sellHistoryGroup := r.Group("/sell-history")
  {
    sellHistoryGroup.GET("", handlers.GetSellHistories)
    sellHistoryGroup.GET("/:id", handlers.GetSellHistoryByID)
    sellHistoryGroup.GET("/card/:card_id", handlers.GetSellHistoriesByCardID)
    sellHistoryGroup.GET("/seller/:seller_id", handlers.GetSellHistoriesBySellerID)
  }

  // Station History routes (read-only)
  stationHistoryGroup := r.Group("/station-history")
  {
    stationHistoryGroup.GET("", handlers.GetStationHistories)
    stationHistoryGroup.GET("/:id", handlers.GetStationHistoryByID)
    stationHistoryGroup.GET("/card/:card_id", handlers.GetStationHistoriesByCardID)
    stationHistoryGroup.GET("/station/:station_id", handlers.GetStationHistoriesByStationID)
    stationHistoryGroup.GET("/action/:action", handlers.GetStationHistoriesByAction)
  }

  // Trip routes
  tripGroup := r.Group("/trip")
  {
    tripGroup.POST("", handlers.CreateTrip)
    tripGroup.GET("", handlers.GetTrips)
    tripGroup.GET("/:id", handlers.GetTripByID)
    tripGroup.PUT("/:id", handlers.UpdateTrip)
    tripGroup.DELETE("/:id", handlers.DeleteTrip)
    tripGroup.GET("/train/:train_id", handlers.GetTripsByTrainID)
    tripGroup.GET("/direction/:direction", handlers.GetTripsByDirection)
    tripGroup.GET("/active", handlers.GetActiveTrips)
  }

  // Train routes
  trainGroup := r.Group("/train")
  {
    trainGroup.POST("", handlers.CreateTrain)
    trainGroup.GET("", handlers.GetTrains)
    trainGroup.GET("/:id", handlers.GetTrainByID)
    trainGroup.PUT("/:id", handlers.UpdateTrain)
    trainGroup.DELETE("/:id", handlers.DeleteTrain)
    trainGroup.GET("/type/:type", handlers.GetTrainsByType)
    trainGroup.GET("/company/:company", handlers.GetTrainsByCompany)
  }

  // Card routes
  cardGroup := r.Group("/card")
  {
    cardGroup.POST("", handlers.CreateCard)                     // Tạo card mới
    cardGroup.GET("", handlers.GetCards)                        // Lấy danh sách tất cả cards
    cardGroup.GET("/:id", handlers.GetCardByID)                 // Lấy card theo ID
    cardGroup.GET("/cardid/:rf_id", handlers.GetCardByCardID)   // Lấy card theo CardID
    cardGroup.PUT("/:rf_id", handlers.UpdateCard)               // Cập nhật card
    cardGroup.DELETE("/:id", handlers.DeleteCard)               // Xóa card
    cardGroup.POST("/:rf_id/topup", handlers.TopUpCard)         // Nạp tiền vào card
    cardGroup.GET("/user/:owner_id", handlers.GetCardsByUser)   // Lấy cards theo owner_id
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
    adminGroup.GET("/users", handlers.GetAllUsers)              // Lấy tất cả users
    adminGroup.GET("/users/simple", handlers.GetAllUsersSimple) // Lấy tất cả users
    adminGroup.GET("/users/statistics", handlers.GetUserStatisticsOptimized)
    adminGroup.GET("/users/:id", handlers.GetUserByID)   // Lấy user theo ID
    adminGroup.PUT("/users/:id", handlers.UpdateUser)    // Cập nhật user
    adminGroup.DELETE("/users/:id", handlers.DeleteUser) // Xóa user
  }

  // Station routes
  stationGroup := r.Group("/station")
  {
    stationGroup.POST("", handlers.CreateStation)         // Tạo station mới
    stationGroup.GET("", handlers.GetStations)            // Lấy danh sách tất cả station
    stationGroup.GET("/:id", handlers.GetStationByID)     // Lấy station theo ID
    stationGroup.PUT("/:id", handlers.UpdateStation)      // Cập nhật station
    stationGroup.DELETE("/:id", handlers.DeleteStation)   // Xóa station
    stationGroup.POST("/:id/checkin", handlers.CheckIn)   // Check-in tại trạm
    stationGroup.POST("/:id/checkout", handlers.CheckOut) // Check-out tại trạm
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
