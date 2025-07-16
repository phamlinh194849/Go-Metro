package handlers

import (
	"go-metro/config"
	"go-metro/models"
	"go-metro/utils"

	"github.com/gin-gonic/gin"
)

// CreateHistory handles POST /history
func CreateHistory(c *gin.Context) {
	var h models.History

	if err := c.ShouldBindJSON(&h); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := config.DB.Create(&h).Error; err != nil {
		utils.InternalServerError(c, "failed to save history")
		return
	}

	utils.SuccessResponse(c, 200, "history saved", h)
}

// GetHistories handles GET /histories
func GetHistories(c *gin.Context) {
	var histories []models.History

	if err := config.DB.Find(&histories).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch histories")
		return
	}

	utils.SuccessResponse(c, 200, "", histories)
}

// GetHistoryByID handles GET /history/:id
func GetHistoryByID(c *gin.Context) {
	id := c.Param("id")
	var history models.History

	if err := config.DB.First(&history, id).Error; err != nil {
		utils.NotFound(c, "history not found")
		return
	}

	utils.SuccessResponse(c, 200, "", history)
} 