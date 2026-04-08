package cheque

// AddLoyaltyDiscounts adds loyalty discounts to the receipt and recalculates the amounts accordingly.
func (cheque *Cheque) AddLoyaltyDiscounts(discounts LoyaltyDiscounts) {
	for _, line := range discounts.ChequeLines.ChequeLine {
		for _, discount := range line.Discounts.Discount {
			cheque.addDiscountToChequeLine(line.ChequeLineNo, discount.DiscountID, discount.Amount)
		}
		pos := cheque.GetLinePosDiscount(line.ChequeLineNo)
		loyalty := cheque.GetLineLoyaltyDiscount(line.ChequeLineNo)
		cheque.recalcReceiptLineAmount(line.ChequeLineNo, pos, loyalty)
	}
	cheque.recalcReceiptTotalAmount()
}

// addDiscountToChequeLine adds a discount to the specified receipt line.
func (cheque *Cheque) addDiscountToChequeLine(chequeLineNo int, id int, amount float64) {
	for i, line := range cheque.ChequeLines.ChequeLines {
		if line.ChequeLineNo == chequeLineNo {
			cheque.ChequeLines.ChequeLines[i].Discounts.Discounts = append(
				cheque.ChequeLines.ChequeLines[i].Discounts.Discounts,
				Discount{
					DiscountID: int32(id),
					Amount:     amount,
				},
			)
		}
	}
}

// GetLinePosDiscount returns the total amount of cashier's discounts applied to the specified receipt line.
func (cheque *Cheque) GetLinePosDiscount(number int) float64 {
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

// GetLineLoyaltyDiscount returns the total amount of loyalty discounts applied to the specified receipt line.
func (cheque *Cheque) GetLineLoyaltyDiscount(number int) float64 {
	var loyalty float64
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

// DeleteLoyaltyDiscounts removes all loyalty discounts from cheque lines, keeping only POS discounts.
func (cheque *Cheque) DeleteLoyaltyDiscounts() {
	cheque.Amount = 0
	for i := range cheque.ChequeLines.ChequeLines {
		pos := Discounts{}
		for _, discount := range cheque.ChequeLines.ChequeLines[i].Discounts.Discounts {
			if discount.DiscountID == 0 {
				pos.Discounts = append(pos.Discounts, discount)
			}
		}
		cheque.ChequeLines.ChequeLines[i].Discounts = pos
		cheque.ChequeLines.ChequeLines[i].Amount = cheque.ChequeLines.ChequeLines[i].Price*
			cheque.ChequeLines.ChequeLines[i].Quantity - cheque.GetLinePosDiscount(cheque.ChequeLines.ChequeLines[i].ChequeLineNo)
		cheque.Amount += cheque.ChequeLines.ChequeLines[i].Amount
	}
}

// recalcReceiptLineAmount recalculates the amount of the specified receipt line.
func (cheque *Cheque) recalcReceiptLineAmount(number int, pos, loyalty float64) {
	for i, line := range cheque.ChequeLines.ChequeLines {
		if line.ChequeLineNo == number {
			cheque.ChequeLines.ChequeLines[i].Amount = line.Quantity*line.Price - pos - loyalty
		}
	}
}

// recalcReceiptTotalAmount recalculates the total amount of the receipt.
func (cheque *Cheque) recalcReceiptTotalAmount() {
	cheque.Amount = 0
	for _, line := range cheque.ChequeLines.ChequeLines {
		cheque.Amount += line.Amount
	}
}
