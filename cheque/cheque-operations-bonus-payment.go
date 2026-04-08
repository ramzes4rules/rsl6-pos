package cheque

// ApplySubtraction applies the specified bonus subtractions to the receipt lines
// and updates the total subtracted bonus and loyalty card information accordingly.
func (cheque *Cheque) ApplySubtraction(subtractions Subtractions, loyaltyCardNumber string) {
	var total float64

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

// GetSubtractedBonus returns the total amount of subtracted bonus from the receipt.
func (cheque *Cheque) GetSubtractedBonus() float64 {
	return cheque.SubtractedBonus
}

// DeleteBonusPayment clears all bonus payment data from the receipt.
func (cheque *Cheque) DeleteBonusPayment() {
	for i := range cheque.ChequeLines.ChequeLines {
		cheque.ChequeLines.ChequeLines[i].BonusDiscount = 0
	}
	cheque.SubtractedBonus = 0
}
