package handlers

import (
  "net/http"
  "strconv"

  "go-metro/config"
  "go-metro/models"
  "go-metro/utils"

  "github.com/gin-gonic/gin"
)

// GetSellHistories handles GET /sell-history
// @Summary Get all sell history records
// @Description Retrieve all sell history records with optional filtering
// @Tags sell-history
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of records per page"
// @Param card_id query string false "Filter by card ID"
// @Param seller_id query string false "Filter by seller ID"
// @Success 200 {object} utils.Response{data=[]models.SellHistory} "Sell histories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /sell-history [get]
func GetSellHistories(c *gin.Context) {
  var sellHistories []models.SellHistory
  query := config.DB.Preload("Card").Preload("Seller")

  // Apply filters
  if cardID := c.Query("card_id"); cardID != "" {
    query = query.Where("card_id = ?", cardID)
  }
  if sellerID := c.Query("seller_id"); sellerID != "" {
    query = query.Where("seller_id = ?", sellerID)
  }

  // Pagination
  page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
  limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
  offset := (page - 1) * limit

  if err := query.Offset(offset).Limit(limit).Find(&sellHistories).Error; err != nil {
    utils.InternalServerError(c, "failed to fetch sell histories")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "sell histories retrieved successfully", sellHistories)
}

// GetSellHistoryByID handles GET /sell-history/:id
// @Summary Get sell history by ID
// @Description Retrieve a specific sell history record by its ID
// @Tags sell-history
// @Accept json
// @Produce json
// @Param id path int true "Sell History ID"
// @Success 200 {object} utils.Response{data=models.SellHistory} "Sell history retrieved successfully"
// @Failure 404 {object} utils.Response "Sell history not found"
// @Router /sell-history/{id} [get]
func GetSellHistoryByID(c *gin.Context) {
  id := c.Param("id")
  var sellHistory models.SellHistory

  if err := config.DB.Preload("Card").Preload("Seller").First(&sellHistory, id).Error; err != nil {
    utils.NotFound(c, "sell history not found")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "sell history retrieved successfully", sellHistory)
}

// GetSellHistoriesByCardID handles GET /sell-history/card/:card_id
// @Summary Get sell histories by card ID
// @Description Retrieve all sell history records for a specific card
// @Tags sell-history
// @Accept json
// @Produce json
// @Param card_id path string true "Card ID"
// @Success 200 {object} utils.Response{data=[]models.SellHistory} "Sell histories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /sell-history/card/{card_id} [get]
func GetSellHistoriesByCardID(c *gin.Context) {
  cardID := c.Param("card_id")
  var sellHistories []models.SellHistory

  if err := config.DB.Preload("Card").Preload("Seller").Where("card_id = ?", cardID).Find(&sellHistories).Error; err != nil {
    utils.InternalServerError(c, "failed to fetch sell histories")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "sell histories retrieved successfully", sellHistories)
}

// GetSellHistoriesBySellerID handles GET /sell-history/seller/:seller_id
// @Summary Get sell histories by seller ID
// @Description Retrieve all sell history records for a specific seller
// @Tags sell-history
// @Accept json
// @Produce json
// @Param seller_id path string true "Seller ID"
// @Success 200 {object} utils.Response{data=[]models.SellHistory} "Sell histories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /sell-history/seller/{seller_id} [get]
func GetSellHistoriesBySellerID(c *gin.Context) {
  sellerID := c.Param("seller_id")
  var sellHistories []models.SellHistory

  if err := config.DB.Preload("Card").Preload("Seller").Where("seller_id = ?", sellerID).Find(&sellHistories).Error; err != nil {
    utils.InternalServerError(c, "failed to fetch sell histories")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "sell histories retrieved successfully", sellHistories)
}
