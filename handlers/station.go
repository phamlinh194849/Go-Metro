package handlers

import (
  "net/http"

  "go-metro/config"
  "go-metro/models"
  "go-metro/utils"

  "github.com/gin-gonic/gin"
)

// StationReq struct for creating station
type StationReq struct {
  Name      string `json:"name" binding:"required"`
  IPAddress string `json:"ip_address"`
}

// CreateStation handles POST /station
// @Summary Create a new station
// @Description Create a new metro station
// @Tags station
// @Accept json
// @Produce json
// @Param station body StationReq true "Station information"
// @Success 200 {object} utils.Response{data=models.Station} "Station created successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station [post]
func CreateStation(c *gin.Context) {
  var stationRequest StationReq

  if err := c.ShouldBindJSON(&stationRequest); err != nil {
    utils.BadRequest(c, err.Error())
    return
  }

  station := models.Station{
    Name:      stationRequest.Name,
    IPAddress: stationRequest.IPAddress,
  }

  if err := config.DB.Create(&station).Error; err != nil {
    utils.InternalServerError(c, "failed to create station")
    return
  }

  utils.SuccessResponse(c, http.StatusCreated, "station created successfully", station)
}

// GetStations handles GET /station
// @Summary Get all stations
// @Description Retrieve all metro stations
// @Tags station
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Station} "Stations retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station [get]
func GetStations(c *gin.Context) {
  var stations []models.Station

  if err := config.DB.Order("id ASC").Find(&stations).Error; err != nil {
    utils.InternalServerError(c, "failed to fetch stations")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "stations retrieved successfully", stations)
}

// GetStationByID handles GET /station/:id
// @Summary Get station by ID
// @Description Retrieve a specific station by its ID
// @Tags station
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} utils.Response{data=models.Station} "Station retrieved successfully"
// @Failure 404 {object} utils.Response "Station not found"
// @Router /station/{id} [get]
func GetStationByID(c *gin.Context) {
  id := c.Param("id")
  var station models.Station

  if err := config.DB.First(&station, id).Error; err != nil {
    utils.NotFound(c, "station not found")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "station retrieved successfully", station)
}

// UpdateStation handles PUT /station/:id
// @Summary Update station
// @Description Update an existing station
// @Tags station
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Param station body StationReq true "Updated station information"
// @Success 200 {object} utils.Response{data=models.Station} "Station updated successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 404 {object} utils.Response "Station not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station/{id} [put]
func UpdateStation(c *gin.Context) {
  id := c.Param("id")
  var station models.Station

  // Check if station exists
  if err := config.DB.First(&station, id).Error; err != nil {
    utils.NotFound(c, "station not found")
    return
  }

  // Bind update data
  var updateData StationReq
  if err := c.ShouldBindJSON(&updateData); err != nil {
    utils.BadRequest(c, err.Error())
    return
  }

  station.Name = updateData.Name
  station.IPAddress = updateData.IPAddress

  if err := config.DB.Save(&station).Error; err != nil {
    utils.InternalServerError(c, "failed to update station")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "station updated successfully", station)
}

// DeleteStation handles DELETE /station/:id
// @Summary Delete station
// @Description Delete a station
// @Tags station
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} utils.Response "Station deleted successfully"
// @Failure 404 {object} utils.Response "Station not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station/{id} [delete]
func DeleteStation(c *gin.Context) {
  id := c.Param("id")
  var station models.Station

  // Check if station exists
  if err := config.DB.First(&station, id).Error; err != nil {
    utils.NotFound(c, "station not found")
    return
  }

  if err := config.DB.Delete(&station).Error; err != nil {
    utils.InternalServerError(c, "failed to delete station")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "station deleted successfully", nil)
}

// CheckInRequest struct for check-in
type CheckInRequest struct {
  CardID string `json:"card_id" binding:"required"`
}

// CheckIn handles POST /station/:id/checkin
// @Summary Check in at station
// @Description Check in a card at a specific station
// @Tags station
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Param request body CheckInRequest true "Check-in information"
// @Success 200 {object} utils.Response{data=models.StationHistory} "Check-in successful"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 404 {object} utils.Response "Station or card not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station/{id}/checkin [post]
func CheckIn(c *gin.Context) {
  stationID := c.Param("id")
  var request CheckInRequest

  if err := c.ShouldBindJSON(&request); err != nil {
    utils.BadRequest(c, err.Error())
    return
  }

  // Check if station exists
  var station models.Station
  if err := config.DB.First(&station, stationID).Error; err != nil {
    utils.NotFound(c, "station not found")
    return
  }

  // Check if card exists
  var card models.Card
  if err := config.DB.Where("rf_id = ?", request.CardID).First(&card).Error; err != nil {
    utils.NotFound(c, "card not found")
    return
  }

  // Check if card has sufficient balance (minimum 5000 VND for check-in)
  if card.Balance < 5000 {
    utils.BadRequest(c, "insufficient balance for check-in")
    return
  }

  // Bắt đầu transaction
  tx := config.DB.Begin()

  // Tạo StationHistory log cho check-in
  if err := utils.CreateStationHistoryLog("checkin", request.CardID, station.ID, 0); err != nil {
    tx.Rollback()
    utils.InternalServerError(c, "failed to create check-in history")
    return
  }

  // Commit transaction
  tx.Commit()

  utils.SuccessResponse(c, http.StatusOK, "check-in successful", gin.H{
    "card_id":    request.CardID,
    "station_id": stationID,
    "action":     "checkin",
    "balance":    card.Balance,
  })
}

// CheckOutRequest struct for check-out
type CheckOutRequest struct {
  CardID string `json:"card_id" binding:"required"`
}

// CheckOut handles POST /station/:id/checkout
// @Summary Check out at station
// @Description Check out a card at a specific station and deduct fare
// @Tags station
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Param request body CheckOutRequest true "Check-out information"
// @Success 200 {object} utils.Response{data=models.StationHistory} "Check-out successful"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 404 {object} utils.Response "Station or card not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /station/{id}/checkout [post]
func CheckOut(c *gin.Context) {
  stationID := c.Param("id")
  var request CheckOutRequest

  if err := c.ShouldBindJSON(&request); err != nil {
    utils.BadRequest(c, err.Error())
    return
  }

  // Check if station exists
  var station models.Station
  if err := config.DB.First(&station, stationID).Error; err != nil {
    utils.NotFound(c, "station not found")
    return
  }

  // Check if card exists
  var card models.Card
  if err := config.DB.Where("rf_id = ?", request.CardID).First(&card).Error; err != nil {
    utils.NotFound(c, "card not found")
    return
  }

  // Calculate fare (fixed fare for demo: 5000 VND)
  fare := 5000.0

  // Check if card has sufficient balance
  if card.Balance < fare {
    utils.BadRequest(c, "insufficient balance for check-out")
    return
  }

  // Bắt đầu transaction
  tx := config.DB.Begin()

  // Deduct fare from card balance
  oldBalance := card.Balance
  card.Balance -= fare
  if err := tx.Save(&card).Error; err != nil {
    tx.Rollback()
    utils.InternalServerError(c, "failed to deduct fare")
    return
  }

  // Tạo StationHistory log cho check-out
  if err := utils.CreateStationHistoryLog("checkout", request.CardID, station.ID, fare); err != nil {
    tx.Rollback()
    utils.InternalServerError(c, "failed to create check-out history")
    return
  }

  // Tạo History log cho payment
  //if err := utils.CreateCardPaymentHistory(request.CardID, card.Username, fare, card.Balance); err != nil {
  //  tx.Rollback()
  //  utils.InternalServerError(c, "failed to create payment history")
  //  return
  //}

  // Commit transaction
  tx.Commit()

  utils.SuccessResponse(c, http.StatusOK, "check-out successful", gin.H{
    "card_id":     request.CardID,
    "station_id":  stationID,
    "action":      "checkout",
    "fare":        fare,
    "old_balance": oldBalance,
    "new_balance": card.Balance,
  })
}
