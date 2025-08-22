package utils

import (
	"go-metro/config"
	"go-metro/consts"
	"go-metro/models"
	"time"
)

// CreateHistoryLog tạo lịch sử giao dịch chung
func CreateHistoryLog(cardID string, userID string, balance float64, userAction consts.UserAction, cardAction consts.CardAction) error {
	history := models.History{
		CardID:     cardID,
		Time:       time.Now(),
		UserID:     userID,
		Balance:    balance,
		UserAction: userAction,
		CardAction: cardAction,
	}

	return config.DB.Create(&history).Error
}

// CreateSellHistoryLog tạo lịch sử bán thẻ
func CreateSellHistoryLog(cardID string, sellerID uint, cardPriceSold float64) error {
	sellHistory := models.SellHistory{
		CardID:        cardID,
		SellerID:      sellerID,
		CardPriceSold: cardPriceSold,
		Time:          time.Now(),
	}

	return config.DB.Create(&sellHistory).Error
}

// CreateStationHistoryLog tạo lịch sử check-in/check-out tại trạm
func CreateStationHistoryLog(action string, cardID string, stationID uint, usedBalance float64) error {
	stationHistory := models.StationHistory{
		Action:      action,
		Time:        time.Now(),
		CardID:      cardID,
		StationID:   stationID,
		UsedBalance: usedBalance,
	}

	return config.DB.Create(&stationHistory).Error
}

// CreateCardTopupHistory tạo lịch sử nạp tiền thẻ
func CreateCardTopupHistory(cardID string, userID string, amount float64, newBalance float64) error {
	// Tạo history log cho topup
	return CreateHistoryLog(cardID, userID, newBalance, consts.UserActionCheckin, consts.CardActionTopup)
}

// CreateCardPaymentHistory tạo lịch sử thanh toán thẻ
func CreateCardPaymentHistory(cardID string, userID string, amount float64, newBalance float64) error {
	// Tạo history log cho payment
	return CreateHistoryLog(cardID, userID, newBalance, consts.UserActionCheckout, consts.CardActionPay)
}

// CreateCardRefundHistory tạo lịch sử hoàn tiền thẻ
func CreateCardRefundHistory(cardID string, userID string, amount float64, newBalance float64) error {
	// Tạo history log cho refund
	return CreateHistoryLog(cardID, userID, newBalance, consts.UserActionCheckin, consts.CardActionRefund)
} 