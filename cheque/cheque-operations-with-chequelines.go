package cheque

import "math"

// NewChequeLine creates a new cheque line with the given parameters.
func NewChequeLine(no int, code string, price, quantity, posd, minp, mina, minb, maxd float64, coupons []string) Line {
	line := Line{
		ChequeLineNo:                   no,
		Price:                          price,
		Quantity:                       quantity,
		Amount:                         price * quantity,
		MinAmount:                      mina,
		MinPrice:                       minp,
		MaxDiscount:                    maxd,
		BonusDiscount:                  0,
		MinAmountAfterCurrencyDiscount: minb,
		Item:                           Item{ItemUID: code},
		Discounts:                      Discounts{},
	}

	if posd > 0 {
		line.AddPosDiscount(posd)
		line.Amount -= posd
	}

	for _, coupon := range coupons {
		line.AddCoupon(coupon)
	}
	return line
}

// AddLineWithChequeLine adds a pre-built cheque line to the receipt.
func (cheque *Cheque) AddLineWithChequeLine(line Line) {
	var coupons []string
	for _, cp := range line.Coupon {
		coupons = append(coupons, cp.CouponNo)
	}
	cheque.AddLine(line.Item.ItemUID, line.Price, line.Quantity, line.MinPrice, line.MinAmount, line.MaxDiscount,
		line.GetPosDiscount(), line.MinAmountAfterCurrencyDiscount, coupons)
}

// AddLine adds a new cheque line to the receipt.
func (cheque *Cheque) AddLine(item string, price, quantity, minPrice, minAmount, maxDiscount, posDiscount, minBonusAmount float64, coupons []string) {
	cheque.DeleteLoyaltyDiscounts()
	cheque.DeleteMessages()

	line := Line{
		ChequeLineNo:                   len(cheque.ChequeLines.ChequeLines) + 1,
		Price:                          price,
		Quantity:                       quantity,
		Amount:                         price * quantity,
		MinPrice:                       minPrice,
		MinAmount:                      minAmount,
		MinAmountAfterCurrencyDiscount: minBonusAmount,
		MaxDiscount:                    maxDiscount,
		Item:                           Item{ItemUID: item},
	}

	for _, coupon := range coupons {
		line.Coupon = append(line.Coupon, Coupon{CouponNo: coupon})
	}

	if posDiscount != 0 {
		line.Discounts.Discounts = append(line.Discounts.Discounts, Discount{
			DiscountID: 0,
			Amount:     posDiscount,
		})
		line.Amount -= posDiscount
	}

	cheque.PositionCount++
	cheque.Amount += line.Amount
	cheque.Amount = math.Round(cheque.Amount*100) / 100
	cheque.ChequeLines.ChequeLines = append(cheque.ChequeLines.ChequeLines, line)
}

// UpdateLine replaces an existing cheque line by its number.
func (cheque *Cheque) UpdateLine(number int, line Line) {
	for i, cl := range cheque.ChequeLines.ChequeLines {
		if cl.ChequeLineNo == number {
			cheque.DeleteLoyaltyDiscounts()
			cheque.DeleteMessages()
			cheque.ChequeLines.ChequeLines[i] = line

			cheque.Amount = 0
			for _, l := range cheque.ChequeLines.ChequeLines {
				cheque.Amount += l.Amount
			}
			break
		}
	}
}

// ChangeLine modifies an existing cheque line at the given index.
func (cheque *Cheque) ChangeLine(number int, item string, price, quantity, minPrice, minAmount, maxDiscount, posDiscount, minBonusAmount float64, coupons []string) {
	cheque.DeleteLoyaltyDiscounts()
	cheque.DeleteMessages()

	cl := &cheque.ChequeLines.ChequeLines[number]
	cl.Item.ItemUID = item
	cl.Price = price
	cl.Quantity = quantity
	cl.Amount = price * quantity
	cl.MinPrice = math.Round(minPrice*100) / 100
	cl.MinAmount = minAmount
	cl.MinAmountAfterCurrencyDiscount = minBonusAmount
	cl.MaxDiscount = maxDiscount
	cl.Coupon = nil
	for _, coupon := range coupons {
		cl.Coupon = append(cl.Coupon, Coupon{CouponNo: coupon})
	}

	cl.Discounts.Discounts = nil
	if posDiscount != 0 {
		cl.Discounts.Discounts = append(cl.Discounts.Discounts, Discount{
			DiscountID: 0,
			Amount:     posDiscount,
		})
		cl.Amount -= posDiscount
	}

	cheque.Amount = 0
	for _, line := range cheque.ChequeLines.ChequeLines {
		cheque.Amount += line.Amount
	}
	cheque.Amount = math.Round(cheque.Amount*100) / 100
}

// DeleteLine removes a cheque line by its number and re-numbers the remaining lines.
func (cheque *Cheque) DeleteLine(index int) {
	var count int
	var remaining []Line

	for _, line := range cheque.ChequeLines.ChequeLines {
		if line.ChequeLineNo != index {
			count++
			line.ChequeLineNo = count
			remaining = append(remaining, line)
		}
	}

	cheque.ChequeLines.ChequeLines = remaining
	cheque.Amount = 0
	for _, line := range cheque.ChequeLines.ChequeLines {
		cheque.Amount += line.Amount
	}
	cheque.PositionCount = int32(len(cheque.ChequeLines.ChequeLines))
}
