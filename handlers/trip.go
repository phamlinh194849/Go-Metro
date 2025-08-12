package handlers

import (
	"go-metro/config"
	"go-metro/models"
	"go-metro/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTrip handles POST /trip
// @Summary Create a new trip record
// @Description Create a new trip record for train journeys
// @Tags trip
// @Accept json
// @Produce json
// @Param trip body models.Trip true "Trip information"
// @Success 200 {object} utils.Response{data=models.Trip} "Trip created successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /trip [post]
func CreateTrip(c *gin.Context) {
	var trip models.Trip

	if err := c.ShouldBindJSON(&trip); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := config.DB.Create(&trip).Error; err != nil {
		utils.InternalServerError(c, "failed to create trip")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "trip created successfully", trip)
}

// GetTrips handles GET /trip
// @Summary Get all trip records
// @Description Retrieve all trip records with optional filtering
// @Tags trip
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of records per page"
// @Param train_id query int false "Filter by train ID"
// @Param direction query string false "Filter by direction"
// @Success 200 {object} utils.Response{data=[]models.Trip} "Trips retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /trip [get]
func GetTrips(c *gin.Context) {
	var trips []models.Trip
	query := config.DB.Preload("Train")

	// Apply filters
	if trainID := c.Query("train_id"); trainID != "" {
		query = query.Where("train_id = ?", trainID)
	}
	if direction := c.Query("direction"); direction != "" {
		query = query.Where("direction = ?", direction)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Order("start_time DESC").Find(&trips).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch trips")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "trips retrieved successfully", trips)
}

// GetTripByID handles GET /trip/:id
// @Summary Get trip by ID
// @Description Retrieve a specific trip record by its ID
// @Tags trip
// @Accept json
// @Produce json
// @Param id path int true "Trip ID"
// @Success 200 {object} utils.Response{data=models.Trip} "Trip retrieved successfully"
// @Failure 404 {object} utils.Response "Trip not found"
// @Router /trip/{id} [get]
func GetTripByID(c *gin.Context) {
	id := c.Param("id")
	var trip models.Trip

	if err := config.DB.Preload("Train").First(&trip, id).Error; err != nil {
		utils.NotFound(c, "trip not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "trip retrieved successfully", trip)
}

// UpdateTrip handles PUT /trip/:id
// @Summary Update trip
// @Description Update an existing trip record
// @Tags trip
// @Accept json
// @Produce json
// @Param id path int true "Trip ID"
// @Param trip body models.Trip true "Updated trip information"
// @Success 200 {object} utils.Response{data=models.Trip} "Trip updated successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 404 {object} utils.Response "Trip not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /trip/{id} [put]
func UpdateTrip(c *gin.Context) {
	id := c.Param("id")
	var trip models.Trip

	// Check if record exists
	if err := config.DB.First(&trip, id).Error; err != nil {
		utils.NotFound(c, "trip not found")
		return
	}

	// Bind update data
	if err := c.ShouldBindJSON(&trip); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := config.DB.Save(&trip).Error; err != nil {
		utils.InternalServerError(c, "failed to update trip")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "trip updated successfully", trip)
}

// DeleteTrip handles DELETE /trip/:id
// @Summary Delete trip
// @Description Delete a trip record
// @Tags trip
// @Accept json
// @Produce json
// @Param id path int true "Trip ID"
// @Success 200 {object} utils.Response "Trip deleted successfully"
// @Failure 404 {object} utils.Response "Trip not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /trip/{id} [delete]
func DeleteTrip(c *gin.Context) {
	id := c.Param("id")
	var trip models.Trip

	// Check if record exists
	if err := config.DB.First(&trip, id).Error; err != nil {
		utils.NotFound(c, "trip not found")
		return
	}

	if err := config.DB.Delete(&trip).Error; err != nil {
		utils.InternalServerError(c, "failed to delete trip")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "trip deleted successfully", nil)
}

// GetTripsByTrainID handles GET /trip/train/:train_id
// @Summary Get trips by train ID
// @Description Retrieve all trip records for a specific train
// @Tags trip
// @Accept json
// @Produce json
// @Param train_id path int true "Train ID"
// @Success 200 {object} utils.Response{data=[]models.Trip} "Trips retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /trip/train/{train_id} [get]
func GetTripsByTrainID(c *gin.Context) {
	trainID := c.Param("train_id")
	var trips []models.Trip

	if err := config.DB.Preload("Train").Where("train_id = ?", trainID).Order("start_time DESC").Find(&trips).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch trips")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "trips retrieved successfully", trips)
}

// GetTripsByDirection handles GET /trip/direction/:direction
// @Summary Get trips by direction
// @Description Retrieve all trip records for a specific direction
// @Tags trip
// @Accept json
// @Produce json
// @Param direction path string true "Direction"
// @Success 200 {object} utils.Response{data=[]models.Trip} "Trips retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /trip/direction/{direction} [get]
func GetTripsByDirection(c *gin.Context) {
	direction := c.Param("direction")
	var trips []models.Trip

	if err := config.DB.Preload("Train").Where("direction = ?", direction).Order("start_time DESC").Find(&trips).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch trips")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "trips retrieved successfully", trips)
}

// GetActiveTrips handles GET /trip/active
// @Summary Get active trips
// @Description Retrieve all currently active trips (where end_time is null or in the future)
// @Tags trip
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Trip} "Active trips retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /trip/active [get]
func GetActiveTrips(c *gin.Context) {
	var trips []models.Trip

	if err := config.DB.Preload("Train").Where("end_time IS NULL OR end_time > NOW()").Order("start_time DESC").Find(&trips).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch active trips")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "active trips retrieved successfully", trips)
} 