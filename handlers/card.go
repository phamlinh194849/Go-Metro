package handlers

import (
	"go-metro/config"
	"go-metro/models"
	"go-metro/utils"

	"github.com/gin-gonic/gin"
)

// CreateCard handles POST /card
// @Summary Create a new card
// @Description Create a new metro card with card ID and user ID
// @Tags card
// @Accept json
// @Produce json
// @Param card body models.Card true "Card information"
// @Success 201 {object} utils.Response{data=models.Card} "Card created successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error or card already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /card [post]
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
// @Summary Get all cards
// @Description Retrieve all metro cards in the system
// @Tags card
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Card} "Cards retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /card [get]
func GetCards(c *gin.Context) {
	var cards []models.Card

	if err := config.DB.Find(&cards).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch cards")
		return
	}

	utils.SuccessResponse(c, 200, "", cards)
}

// GetCardByID handles GET /card/:id
// @Summary Get card by ID
// @Description Retrieve a specific card by its database ID
// @Tags card
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} utils.Response{data=models.Card} "Card retrieved successfully"
// @Failure 404 {object} utils.Response "Card not found"
// @Router /card/{id} [get]
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
// @Summary Get card by Card ID
// @Description Retrieve a specific card by its card ID (physical card number)
// @Tags card
// @Accept json
// @Produce json
// @Param card_id path string true "Card ID (physical card number)"
// @Success 200 {object} utils.Response{data=models.Card} "Card retrieved successfully"
// @Failure 404 {object} utils.Response "Card not found"
// @Router /card/cardid/{card_id} [get]
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
// @Summary Update card
// @Description Update a specific card's information
// @Tags card
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param card body models.Card true "Updated card information"
// @Success 200 {object} utils.Response{data=models.Card} "Card updated successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 404 {object} utils.Response "Card not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /card/{id} [put]
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
// @Summary Delete card
// @Description Delete a specific card from the system
// @Tags card
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} utils.Response "Card deleted successfully"
// @Failure 404 {object} utils.Response "Card not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /card/{id} [delete]
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
// @Summary Top up card balance
// @Description Add money to a card's balance
// @Tags card
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
//
//	@Param request body object true "Top-up amount" schema={
//	  "type": "object",
//	  "required": ["amount"],
//	  "properties": {
//	    "amount": {"type": "number", "minimum": 0.01}
//	  }
//	}
//
// @Success 200 {object} utils.Response{data=models.Card} "Card topped up successfully"
// @Failure 400 {object} utils.Response "Bad request - invalid amount"
// @Failure 404 {object} utils.Response "Card not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /card/{id}/topup [post]
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
// @Summary Get cards by user ID
// @Description Retrieve all cards belonging to a specific user
// @Tags card
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} utils.Response{data=[]models.Card} "User cards retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /card/user/{user_id} [get]
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
// @Summary Get cards by status
// @Description Retrieve all cards with a specific status (active, inactive, blocked)
// @Tags card
// @Accept json
// @Produce json
// @Param status path string true "Card status" Enums(active, inactive, blocked)
// @Success 200 {object} utils.Response{data=[]models.Card} "Cards by status retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /card/status/{status} [get]
func GetCardsByStatus(c *gin.Context) {
	status := c.Param("status")
	var cards []models.Card

	if err := config.DB.Where("status = ?", status).Find(&cards).Error; err != nil {
		utils.InternalServerError(c, "failed to fetch cards by status")
		return
	}

	utils.SuccessResponse(c, 200, "", cards)
}
