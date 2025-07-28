package handlers

import (
	"fmt"
	"go-metro/config"
	"go-metro/consts"
	"go-metro/models"
	"go-metro/utils"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// CardRequest struct for creating card (without rf_id as it's auto-generated)
type CardReq struct {
	OwnerID string `json:"owner_id"`
	Type    string `json:"type"  binding:"required"`
}

type UpdateCardReq struct {
	OwnerID string  `json:"owner_id"`
	Balance float64 `json:"balance" gorm:"default:0"`
	Status  string  `json:"status"`
	Type    int
}

func generateCardID() string {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Int63n(10000000000) // Generate 10-digit number
	return fmt.Sprintf("GM%010d", randomNumber)
}

func isCardIDUnique(cardID string) bool {
	var count int64
	config.DB.Model(&models.Card{}).Where("rf_id = ?", cardID).Count(&count)
	return count == 0
}

func generateUniqueCardID() string {
	var cardID string
	for {
		cardID = generateCardID()
		if isCardIDUnique(cardID) {
			break
		}
	}
	return cardID
}

// OK
// CreateCard handles POST /card
// @Summary Create a new card
// @Description Create a new metro card with auto-generated card ID and user ID
// @Tags card
// @Accept json
// @Produce json
// @Param card body CardReq true "Card information"
// @Router /card [post]
func CreateCard(c *gin.Context) {
	var cardRequest CardReq

	if err := c.ShouldBindJSON(&cardRequest); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	cardID := generateUniqueCardID()
	var user models.User
	if err := config.DB.Where("username = ?", cardRequest.OwnerID).First(&user).Error; err != nil {
		utils.NotFound(c, "Người dùng không tồn tại")
		return
	}

	card := models.Card{
		RFID:     cardID,
		Username: cardRequest.OwnerID,
	}

	// Search OwnerID trong database trước

	if card.Username != "" {
		card.Status = consts.ActiveStatus
	} else {
		card.Status = consts.InactiveStatus
	}

	if cardRequest.Type == consts.StudentCard.ToText() {
		card.Balance = consts.StudentCard.ToDefaultBlance()
		card.Price = consts.StudentCard.ToPrice()
		card.Type = consts.StudentCard
	} else if cardRequest.Type == consts.NormalCard.ToText() {
		card.Balance = consts.NormalCard.ToDefaultBlance()
		card.Price = consts.NormalCard.ToPrice()
		card.Type = consts.NormalCard
	} else if cardRequest.Type == consts.VipCard.ToText() {
		card.Balance = consts.VipCard.ToDefaultBlance()
		card.Price = consts.VipCard.ToPrice()
		card.Type = consts.VipCard
	}

	if err := config.DB.Create(&card).Error; err != nil {
		utils.InternalServerError(c, "Lỗi tạo thẻ")
		return
	}

	utils.SuccessResponse(c, 201, "Tạo thẻ thành công", card)
}

// OK
// GetCards handles GET /card
// @Summary Get all cards
// @Description Retrieve all metro cards in the system
// @Tags card
// @Accept json
// @Produce json
// @Router /card [get]
func GetCards(c *gin.Context) {
	var cards []models.Card

	if err := config.DB.Find(&cards).Error; err != nil {
		utils.InternalServerError(c, "không lấy được danh sách thẻ")
		return
	}

	utils.SuccessResponse(c, 200, "", cards)
}

// OK
// GetCardByID handles GET /card/:id
// @Summary Get card by ID
// @Description Retrieve a specific card by its database ID
// @Tags card
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} utils.Response{data=models.Card} "Card retrieved successfully"
// @Failure 404 {object} utils.Response "Thẻ không tồn tại"
// @Router /card/{id} [get]
func GetCardByID(c *gin.Context) {
	id := c.Param("id")
	var card models.Card

	if err := config.DB.First(&card, id).Error; err != nil {
		utils.NotFound(c, "Thẻ không tồn tại")
		return
	}

	utils.SuccessResponse(c, 200, "", card)
}

// OK
// GetCardByCardID handles GET /card/cardid/:rf_id
// @Summary Get card by Card ID
// @Description Retrieve a specific card by its card ID (physical card number)
// @Tags card
// @Accept json
// @Produce json
// @Param rf_id path string true "Card ID (physical card number)"
// @Router /card/cardid/{rf_id} [get]
func GetCardByCardID(c *gin.Context) {
	cardID := c.Param("rf_id")
	var card models.Card

	if err := config.DB.Where("rf_id = ?", cardID).First(&card).Error; err != nil {
		utils.NotFound(c, "Thẻ không tồn tại")
		return
	}

	utils.SuccessResponse(c, 200, "", card)
}

// TODO
// UpdateCard handles PUT /card/:rf_id
// @Summary Update card
// @Description Update a specific card's information
// @Tags card
// @Accept json
// @Produce json
// @Param rf_id path string true "Card ID"
// @Param card body UpdateCardReq true "Updated card information"
// @Router /card/{rf_id} [put]
func UpdateCard(c *gin.Context) {
	rf_id := c.Param("rf_id")
	var card models.Card

	// Check if card exists
	if err := config.DB.Where("rf_id = ?", rf_id).First(&card).Error; err != nil {
		utils.NotFound(c, "Thẻ không tồn tại")
		return
	}

	// Bind update data
	var updateData UpdateCardReq
	updateDataCard := models.Card{
		Username: updateData.OwnerID,
		Balance:  updateData.Balance,
		Status:   consts.Status(updateData.Status),
		Type:     consts.CardType(updateData.Type),
	}

	switch updateDataCard.Type {
	case consts.StudentCard:
		updateDataCard.Price = consts.StudentCard.ToPrice()
	case consts.NormalCard:
		updateDataCard.Price = consts.NormalCard.ToPrice()
	case consts.VipCard:
		updateDataCard.Price = consts.VipCard.ToPrice()
	}

	if err := c.ShouldBindJSON(&updateDataCard); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Update card
	if err := config.DB.Model(&card).Updates(updateData).Error; err != nil {
		utils.InternalServerError(c, "Lỗi cập nhật thẻ")
		return
	}

	// Get updated card
	config.DB.First(&card, rf_id)
	utils.SuccessResponse(c, 200, "Cập nhật thành công", card)
}

// DeleteCard handles DELETE /card/:id
// @Summary Delete card
// @Description Delete a specific card from the system
// @Tags card
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Router /card/{id} [delete]
func DeleteCard(c *gin.Context) {
	id := c.Param("id")
	var card models.Card

	if err := config.DB.First(&card, id).Error; err != nil {
		utils.NotFound(c, "Thẻ không tồn tại")
		return
	}

	if err := config.DB.Delete(&card).Error; err != nil {
		utils.InternalServerError(c, "failed to delete card")
		return
	}

	utils.SuccessResponse(c, 200, "card deleted successfully", nil)
}

// TopUpCard handles POST /card/:rf_id/topup
// @Summary Top up card balance
// @Description Add money to a card's balance
// @Tags card
// @Accept json
// @Produce json
// @Param rf_id path string true "Card ID"
// @Param request body object true "Top-up amount"
// @Router /card/{rf_id}/topup [post]
func TopUpCard(c *gin.Context) {
	id := c.Param("rf_id")
	var card models.Card

	// Check if card exists
	if err := config.DB.Where("rf_id = ?", id).First(&card).Error; err != nil {
		utils.NotFound(c, "Thẻ không tồn tại")
		return
	}

	// Parse amount from request
	var request struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, "số tiền không hợp lệ")
		return
	}

	// Update balance
	card.Balance += request.Amount
	if err := config.DB.Save(&card).Error; err != nil {
		utils.InternalServerError(c, "Nạp tiền thất bại")
		return
	}

	utils.SuccessResponse(c, 200, "Nạp tiền thành công", card)
}

// GetCardsByUser handles GET /card/user/:owner_id
// @Summary Get cards by user ID
// @Description Retrieve all cards belonging to a specific user
// @Tags card
// @Accept json
// @Produce json
// @Param owner_id path string true "User ID"
// @Router /card/user/{owner_id} [get]
func GetCardsByUser(c *gin.Context) {
	username := c.Param("owner_id")
	var cards []models.Card

	if err := config.DB.Where("username = ?", username).Find(&cards).Error; err != nil {
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
