package handlers

import (
  "net/http"
  "strconv"

  "go-metro/config"
  "go-metro/models"
  "go-metro/utils"

  "github.com/gin-gonic/gin"
)

// GetStationHistories handles GET /station-history
// @Summary Get all station history records
// @Description Retrieve all station history records with optional filtering
// @Tags station-history
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of records per page"
// @Param card_id query string false "Filter by card ID"
// @Param station_id query int false "Filter by station ID"
// @Param action query string false "Filter by action (checkin/checkout)"
// @Success 200 {object} utils.Response{data=[]models.StationHistory} "Station histories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station-history [get]
func GetStationHistories(c *gin.Context) {
  var stationHistories []models.StationHistory
  query := config.DB.Preload("Card").Preload("Station")

  // Apply filters
  if cardID := c.Query("card_id"); cardID != "" {
    query = query.Where("card_id = ?", cardID)
  }
  if stationID := c.Query("station_id"); stationID != "" {
    query = query.Where("station_id = ?", stationID)
  }
  if action := c.Query("action"); action != "" {
    query = query.Where("action = ?", action)
  }

  // Pagination
  page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
  limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
  offset := (page - 1) * limit

  if err := query.Offset(offset).Limit(limit).Order("time DESC").Find(&stationHistories).Error; err != nil {
    utils.InternalServerError(c, "failed to fetch station histories")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "station histories retrieved successfully", stationHistories)
}

// GetStationHistoryByID handles GET /station-history/:id
// @Summary Get station history by ID
// @Description Retrieve a specific station history record by its ID
// @Tags station-history
// @Accept json
// @Produce json
// @Param id path int true "Station History ID"
// @Success 200 {object} utils.Response{data=models.StationHistory} "Station history retrieved successfully"
// @Failure 404 {object} utils.Response "Station history not found"
// @Router /station-history/{id} [get]
func GetStationHistoryByID(c *gin.Context) {
  id := c.Param("id")
  var stationHistory models.StationHistory

  if err := config.DB.Preload("Card").Preload("Station").First(&stationHistory, id).Error; err != nil {
    utils.NotFound(c, "station history not found")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "station history retrieved successfully", stationHistory)
}

// GetStationHistoriesByCardID handles GET /station-history/card/:card_id
// @Summary Get station histories by card ID
// @Description Retrieve all station history records for a specific card
// @Tags station-history
// @Accept json
// @Produce json
// @Param card_id path string true "Card ID"
// @Success 200 {object} utils.Response{data=[]models.StationHistory} "Station histories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station-history/card/{card_id} [get]
func GetStationHistoriesByCardID(c *gin.Context) {
  cardID := c.Param("card_id")
  var stationHistories []models.StationHistory

  if err := config.DB.Preload("Card").Preload("Station").Where("card_id = ?", cardID).Order("time DESC").Find(&stationHistories).Error; err != nil {
    utils.InternalServerError(c, "failed to fetch station histories")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "station histories retrieved successfully", stationHistories)
}

// GetStationHistoriesByStationID handles GET /station-history/station/:station_id
// @Summary Get station histories by station ID
// @Description Retrieve all station history records for a specific station
// @Tags station-history
// @Accept json
// @Produce json
// @Param station_id path int true "Station ID"
// @Success 200 {object} utils.Response{data=[]models.StationHistory} "Station histories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station-history/station/{station_id} [get]
func GetStationHistoriesByStationID(c *gin.Context) {
  stationID := c.Param("station_id")
  var stationHistories []models.StationHistory

  if err := config.DB.Preload("Card").Preload("Station").Where("station_id = ?", stationID).Order("time DESC").Find(&stationHistories).Error; err != nil {
    utils.InternalServerError(c, "failed to fetch station histories")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "station histories retrieved successfully", stationHistories)
}

// GetStationHistoriesByAction handles GET /station-history/action/:action
// @Summary Get station histories by action
// @Description Retrieve all station history records for a specific action (checkin/checkout)
// @Tags station-history
// @Accept json
// @Produce json
// @Param action path string true "Action (checkin/checkout)"
// @Success 200 {object} utils.Response{data=[]models.StationHistory} "Station histories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station-history/action/{action} [get]
func GetStationHistoriesByAction(c *gin.Context) {
  action := c.Param("action")
  var stationHistories []models.StationHistory

  if err := config.DB.Preload("Card").Preload("Station").Where("action = ?", action).Order("time DESC").Find(&stationHistories).Error; err != nil {
    utils.InternalServerError(c, "failed to fetch station histories")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "station histories retrieved successfully", stationHistories)
}
