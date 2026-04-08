package cheque

import "github.com/ramzes4rules/rsl6-pos/models"

// ApplyGiftCardPayment applies a gift card payment to the receipt
// by adding card entry with the specified number and amount
func (cheque *models.Cheque) ApplyGiftCardPayment(number string, amount float64) {
	cheque.DiscountCard = append(cheque.DiscountCard, models.DiscountCard{
		DiscountCardNo:       number,
		SubtractAmount:       amount,
		BonusCard:            false,
		EnteredAsPhoneNumber: false,
	})
}
