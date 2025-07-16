package main

import (
	"go-metro/config"
	"go-metro/models"
	"go-metro/routes"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.ConnectDB()

	// Kiểm tra biến MIGRATE trong env
	migrate := strings.ToLower(os.Getenv("MIGRATE")) == "true"
	if migrate {
		models.MigrateHistory()
		models.MigrateCard()
		models.MigrateUser()
	}

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	r.Run(":8080")
}
