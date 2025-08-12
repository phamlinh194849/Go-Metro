package handlers

import (
	"go-metro/config"
	"go-metro/models"
	"go-metro/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTrain handles POST /train
// @Summary Create a new train record
// @Description Create a new train record
// @Tags train
// @Accept json
// @Produce json
// @Param train body models.Train true "Train information"
// @Success 200 {object} utils.Response{data=models.Train} "Train created successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /train [post]
func CreateTrain(c *gin.Context) {
	var train models.Train

	if err := c.ShouldBindJSON(&train); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := config.DB.Create(&train).Error; err != nil {
		utils.InternalServerError(c, "failed to create train")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "train created successfully", train)
}

// GetTrains handles GET /train
// @Summary Get all train records
// @Description Retrieve all train records with optional filtering
// @Tags train
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of records per page"
// @Param type query string false "Filter by train type"
// @Param company query string false "Filter by company"
// @Success 200 {object} utils.Response{data=[]models.Train} "Trains retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /train [get]
func GetTrains(c *gin.Context) {
	var trains []models.Train
	query := config.DB

	// Apply filters
	if trainType := c.Query("type"); trainType != "" {
		query = query.Where("type = ?", trainType)
	}
	if company := c.Query("company"); company != "" {
		query = query.Where("company = ?", company)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&trains).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch trains")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "trains retrieved successfully", trains)
}

// GetTrainByID handles GET /train/:id
// @Summary Get train by ID
// @Description Retrieve a specific train record by its ID
// @Tags train
// @Accept json
// @Produce json
// @Param id path int true "Train ID"
// @Success 200 {object} utils.Response{data=models.Train} "Train retrieved successfully"
// @Failure 404 {object} utils.Response "Train not found"
// @Router /train/{id} [get]
func GetTrainByID(c *gin.Context) {
	id := c.Param("id")
	var train models.Train

	if err := config.DB.First(&train, id).Error; err != nil {
		utils.NotFound(c, "train not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "train retrieved successfully", train)
}

// UpdateTrain handles PUT /train/:id
// @Summary Update train
// @Description Update an existing train record
// @Tags train
// @Accept json
// @Produce json
// @Param id path int true "Train ID"
// @Param train body models.Train true "Updated train information"
// @Success 200 {object} utils.Response{data=models.Train} "Train updated successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 404 {object} utils.Response "Train not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /train/{id} [put]
func UpdateTrain(c *gin.Context) {
	id := c.Param("id")
	var train models.Train

	// Check if record exists
	if err := config.DB.First(&train, id).Error; err != nil {
		utils.NotFound(c, "train not found")
		return
	}

	// Bind update data
	if err := c.ShouldBindJSON(&train); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := config.DB.Save(&train).Error; err != nil {
		utils.InternalServerError(c, "failed to update train")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "train updated successfully", train)
}

// DeleteTrain handles DELETE /train/:id
// @Summary Delete train
// @Description Delete a train record
// @Tags train
// @Accept json
// @Produce json
// @Param id path int true "Train ID"
// @Success 200 {object} utils.Response "Train deleted successfully"
// @Failure 404 {object} utils.Response "Train not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /train/{id} [delete]
func DeleteTrain(c *gin.Context) {
	id := c.Param("id")
	var train models.Train

	// Check if record exists
	if err := config.DB.First(&train, id).Error; err != nil {
		utils.NotFound(c, "train not found")
		return
	}

	if err := config.DB.Delete(&train).Error; err != nil {
		utils.InternalServerError(c, "failed to delete train")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "train deleted successfully", nil)
}

// GetTrainsByType handles GET /train/type/:type
// @Summary Get trains by type
// @Description Retrieve all train records for a specific type
// @Tags train
// @Accept json
// @Produce json
// @Param type path string true "Train type"
// @Success 200 {object} utils.Response{data=[]models.Train} "Trains retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /train/type/{type} [get]
func GetTrainsByType(c *gin.Context) {
	trainType := c.Param("type")
	var trains []models.Train

	if err := config.DB.Where("type = ?", trainType).Find(&trains).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch trains")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "trains retrieved successfully", trains)
}

// GetTrainsByCompany handles GET /train/company/:company
// @Summary Get trains by company
// @Description Retrieve all train records for a specific company
// @Tags train
// @Accept json
// @Produce json
// @Param company path string true "Company name"
// @Success 200 {object} utils.Response{data=[]models.Train} "Trains retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /train/company/{company} [get]
func GetTrainsByCompany(c *gin.Context) {
	company := c.Param("company")
	var trains []models.Train

	if err := config.DB.Where("company = ?", company).Find(&trains).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch trains")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "trains retrieved successfully", trains)
} 