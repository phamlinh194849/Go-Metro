package handlers

import (
  "go-metro/config"
  "go-metro/models"
  "go-metro/utils"

  "github.com/gin-gonic/gin"
)

// GetHistories handles GET /histories
// @Summary Get all history records
// @Description Retrieve all transaction history records
// @Tags history
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.History} "Histories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /history [get]
func GetHistories(c *gin.Context) {
  var histories []models.History

  if err := config.DB.Find(&histories).Error; err != nil {
    utils.InternalServerError(c, "failed to fetch histories")
    return
  }

  utils.SuccessResponse(c, 200, "", histories)
}

// GetHistoryByID handles GET /history/:id
// @Summary Get history by ID
// @Description Retrieve a specific history record by its ID
// @Tags history
// @Accept json
// @Produce json
// @Param id path int true "History ID"
// @Success 200 {object} utils.Response{data=models.History} "History retrieved successfully"
// @Failure 404 {object} utils.Response "History not found"
// @Router /history/{id} [get]
func GetHistoryByID(c *gin.Context) {
  id := c.Param("id")
  var history models.History

  if err := config.DB.First(&history, id).Error; err != nil {
    utils.NotFound(c, "history not found")
    return
  }

  utils.SuccessResponse(c, 200, "", history)
}
