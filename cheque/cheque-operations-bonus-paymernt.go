package cheque

import "github.com/ramzes4rules/rsl6-pos/models"

// ApplySubtraction applies the specified subtractions to the receipt lines and
// updates the total subtracted bonus and loyalty card information accordingly
func (cheque *models.Cheque) ApplySubtraction(subtractions models.Subtractions, loyaltyCardNumber string) {
	var total float64 = 0

	for _, subtraction := range subtractions {
		for i, line := range cheque.ChequeLines.ChequeLines {
			if line.ChequeLineNo == subtraction.ChequeLineNo {
				cheque.ChequeLines.ChequeLines[i].BonusDiscount = subtraction.Amount
				total += subtraction.Amount
			}
		}
	}
	cheque.SubtractedBonus = total

	for i, card := range cheque.DiscountCard {
		if card.DiscountCardNo == loyaltyCardNumber {
			cheque.DiscountCard[i].SubtractedBonus = total
		}
	}
}
func (cheque *models.Cheque) GetSubtractedBonus() float64 {
	var amount float64 = 0
	amount = cheque.SubtractedBonus
	return amount
}
func (cheque *models.Cheque) DeleteBonusPayment() {
	for i, _ := range cheque.ChequeLines.ChequeLines {
		cheque.ChequeLines.ChequeLines[i].BonusDiscount = 0
	}
	cheque.SubtractedBonus = 0
}
