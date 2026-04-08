package cheque

import "github.com/ramzes4rules/rsl6-pos/models"

// AddLoyaltyDiscounts adds loyalty discounts to the receipt and recalculates the amounts accordingly
func (cheque *models.Cheque) AddLoyaltyDiscounts(discounts models.LoyaltyDiscounts) {
	for _, line := range discounts.ChequeLines.ChequeLine {
		for _, discount := range line.Discounts.Discount {
			cheque.AddDiscountToChequeLine(line.ChequeLineNo, discount)
		}
		pos := cheque.GetLinePosDiscount(line.ChequeLineNo)
		loyalty := cheque.GetLineLoyaltyDiscount(line.ChequeLineNo)
		cheque.recalcReceiptLineAmount(line.ChequeLineNo, pos, loyalty)

	}
	cheque.recalcReceiptTotalAmount()
}

// AddDiscountToChequeLine adds a discount to the specified receipt line
func (cheque *models.Cheque) AddDiscountToChequeLine(chequeLineNo int, discount IncomeDiscount) {
	for i, line := range cheque.ChequeLines.ChequeLines {
		if line.ChequeLineNo == chequeLineNo {
			cheque.ChequeLines.ChequeLines[i].Discounts.Discounts = append(cheque.ChequeLines.ChequeLines[i].Discounts.Discounts,
				models.Discount{
					DiscountID: int32(discount.DiscountID),
					Amount:     discount.Amount,
				},
			)
		}
	}
}

// GetLinePosDiscount returns the total amount of cashier's discounts applied to the specified receipt line
func (cheque *models.Cheque) GetLinePosDiscount(number int) float64 {
	var pos float64
	for _, line := range cheque.ChequeLines.ChequeLines {
		if line.ChequeLineNo == number {
			for _, discount := range line.Discounts.Discounts {
				if discount.DiscountID == 0 {
					pos += discount.Amount
				}
			}
		}
	}
	return pos
}

// GetLineLoyaltyDiscount returns the total amount of loyalty discounts applied to the specified receipt line
func (cheque *models.Cheque) GetLineLoyaltyDiscount(number int) float64 {
	var loyalty float64 = 0
	for _, line := range cheque.ChequeLines.ChequeLines {
		if line.ChequeLineNo == number {
			for _, discount := range line.Discounts.Discounts {
				if discount.DiscountID != 0 {
					loyalty += discount.Amount
				}
			}
		}
	}
	return loyalty
}

// DeleteLoyaltyDiscounts removes all discounts record from check lines
func (cheque *models.Cheque) DeleteLoyaltyDiscounts() {
	//log.Trace().Msg("Clearing receipt discounts")
	//old := cheque.Amount
	cheque.Amount = 0
	for i := range cheque.ChequeLines.ChequeLines {
		pos := models.Discounts{}
		for _, discount := range cheque.ChequeLines.ChequeLines[i].Discounts.Discounts {
			if discount.DiscountID == 0 {
				pos.Discounts = append(pos.Discounts, discount)
			}
		}
		cheque.ChequeLines.ChequeLines[i].Discounts = pos
		//log.Trace().Msg("POS discount for")

		//old := cheque.ChequeLines.ChequeLines[i].Amount
		cheque.ChequeLines.ChequeLines[i].Amount = cheque.ChequeLines.ChequeLines[i].Price*
			cheque.ChequeLines.ChequeLines[i].Quantity - cheque.GetLinePosDiscount(cheque.ChequeLines.ChequeLines[i].ChequeLineNo)
		//log.Trace().Int("line", i+1).Float64("before", old).
		//	Float64("after", cheque.ChequeLines.ChequeLines[i].Amount).Msg("Processing line:")
		cheque.Amount += cheque.ChequeLines.ChequeLines[i].Amount
	}
	//log.Trace().Float64("before", old).Float64("after", cheque.Amount).
	//	Msg("Receipt discounts were cleared, receipt amount:")
}

// recalcReceiptLineAmount recalculates the amount of the specified receipt line
// based on the original price, quantity, and applied discounts
func (cheque *models.Cheque) recalcReceiptLineAmount(number int, pos, loyalty float64) {
	for i, line := range cheque.ChequeLines.ChequeLines {
		if line.ChequeLineNo == number {
			cheque.ChequeLines.ChequeLines[i].Amount = cheque.ChequeLines.ChequeLines[i].Quantity*cheque.ChequeLines.ChequeLines[i].Price - pos - loyalty
		}
	}
}

// recalcReceiptTotalAmount recalculates the total amount of the receipt based on the amounts of all receipt lines
func (cheque *models.Cheque) recalcReceiptTotalAmount() {
	cheque.Amount = 0
	for _, line := range cheque.ChequeLines.ChequeLines {
		cheque.Amount += line.Amount
	}
}
