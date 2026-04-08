package cheque

// ApplyGiftCardPayment applies a gift card payment to the receipt
// by adding a card entry with the specified number and amount.
func (cheque *Cheque) ApplyGiftCardPayment(number string, amount float64) {
	cheque.DiscountCard = append(cheque.DiscountCard, DiscountCard{
		DiscountCardNo:       number,
		SubtractAmount:       amount,
		BonusCard:            false,
		EnteredAsPhoneNumber: false,
	})
}
