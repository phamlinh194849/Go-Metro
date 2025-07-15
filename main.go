package main

import (
	"go-metro/db"
	"go-metro/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()
	models.MigrateHistory()

	r := gin.Default()

	r.POST("/history", func(c *gin.Context) {
		var h models.History

		if err := c.ShouldBindJSON(&h); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.DB.Create(&h).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save history"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "history saved", "data": h})
	})

	r.Run(":8080")
}
