package handlers

import (
	"go-metro/config"
	"go-metro/models"
	"go-metro/utils"

	"github.com/gin-gonic/gin"
)

// CreateCard handles POST /card
func CreateCard(c *gin.Context) {
	var card models.Card

	if err := c.ShouldBindJSON(&card); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Check if card already exists
	var existingCard models.Card
	if err := config.DB.Where("card_id = ?", card.CardID).First(&existingCard).Error; err == nil {
		utils.BadRequest(c, "card already exists")
		return
	}

	if err := config.DB.Create(&card).Error; err != nil {
		utils.InternalServerError(c, "failed to create card")
		return
	}

	utils.SuccessResponse(c, 201, "card created successfully", card)
}

// GetCards handles GET /card
func GetCards(c *gin.Context) {
	var cards []models.Card

	if err := config.DB.Find(&cards).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch cards")
		return
	}

	utils.SuccessResponse(c, 200, "", cards)
}

// GetCardByID handles GET /card/:id
func GetCardByID(c *gin.Context) {
	id := c.Param("id")
	var card models.Card

	if err := config.DB.First(&card, id).Error; err != nil {
		utils.NotFound(c, "card not found")
		return
	}

	utils.SuccessResponse(c, 200, "", card)
}

// GetCardByCardID handles GET /card/cardid/:card_id
func GetCardByCardID(c *gin.Context) {
	cardID := c.Param("card_id")
	var card models.Card

	if err := config.DB.Where("card_id = ?", cardID).First(&card).Error; err != nil {
		utils.NotFound(c, "card not found")
		return
	}

	utils.SuccessResponse(c, 200, "", card)
}

// UpdateCard handles PUT /card/:id
func UpdateCard(c *gin.Context) {
	id := c.Param("id")
	var card models.Card

	// Check if card exists
	if err := config.DB.First(&card, id).Error; err != nil {
		utils.NotFound(c, "card not found")
		return
	}

	// Bind update data
	var updateData models.Card
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Update card
	if err := config.DB.Model(&card).Updates(updateData).Error; err != nil {
		utils.InternalServerError(c, "failed to update card")
		return
	}

	// Get updated card
	config.DB.First(&card, id)
	utils.SuccessResponse(c, 200, "card updated successfully", card)
}

// DeleteCard handles DELETE /card/:id
func DeleteCard(c *gin.Context) {
	id := c.Param("id")
	var card models.Card

	if err := config.DB.First(&card, id).Error; err != nil {
		utils.NotFound(c, "card not found")
		return
	}

	if err := config.DB.Delete(&card).Error; err != nil {
		utils.InternalServerError(c, "failed to delete card")
		return
	}

	utils.SuccessResponse(c, 200, "card deleted successfully", nil)
}

// TopUpCard handles POST /card/:id/topup
func TopUpCard(c *gin.Context) {
	id := c.Param("id")
	var card models.Card

	// Check if card exists
	if err := config.DB.First(&card, id).Error; err != nil {
		utils.NotFound(c, "card not found")
		return
	}

	// Parse amount from request
	var request struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, "invalid amount")
		return
	}

	// Update balance
	card.Balance += request.Amount
	if err := config.DB.Save(&card).Error; err != nil {
		utils.InternalServerError(c, "failed to top up card")
		return
	}

	utils.SuccessResponse(c, 200, "card topped up successfully", card)
}

// GetCardsByUser handles GET /card/user/:user_id
func GetCardsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	var cards []models.Card

	if err := config.DB.Where("user_id = ?", userID).Find(&cards).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch user cards")
		return
	}

	utils.SuccessResponse(c, 200, "", cards)
}

// GetCardsByStatus handles GET /card/status/:status
func GetCardsByStatus(c *gin.Context) {
	status := c.Param("status")
	var cards []models.Card

	if err := config.DB.Where("status = ?", status).Find(&cards).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch cards by status")
		return
	}

	utils.SuccessResponse(c, 200, "", cards)
}
