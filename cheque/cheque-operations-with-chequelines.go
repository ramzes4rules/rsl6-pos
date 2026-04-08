package cheque

import (
	"math"

	"github.com/ramzes4rules/rsl6-pos/models"
)

func NewChequeLine(no int, code string, price, quantity, posd, minp, mina, minb, maxd float64, coupons []string) models.Line {

	// Create new cheque line
	line := models.Line{
		ChequeLineNo:                   no,
		Price:                          price,
		Quantity:                       quantity,
		Amount:                         price * quantity,
		MinAmount:                      mina,
		MinPrice:                       minp,
		MaxDiscount:                    maxd,
		BonusDiscount:                  0,
		MinAmountAfterCurrencyDiscount: minb,
		Item:                           models.Item{ItemUID: code},
		Discounts:                      models.Discounts{},
		Coupon:                         nil,
	}

	// Add pos discount
	if posd > 0 {
		line.AddPosDiscount(posd)
		line.Amount -= posd
	}

	// Add coupons
	for _, coupon := range coupons {
		line.AddCoupon(coupon)
	}
	return line
}

func (cheque *models.Cheque) AddLineWithChequeLine(line models.Line) {
	var coupons []string
	for _, cp := range line.Coupon {
		coupons = append(coupons, cp.CouponNo)
	}
	cheque.AddLine(line.Item.ItemUID, line.Price, line.Quantity, line.MinPrice, line.MinAmount, line.MaxDiscount,
		line.GetPosDiscount(), line.MinAmountAfterCurrencyDiscount, coupons)
}

func (cheque *models.Cheque) AddLine(item string, price, quantity, minPrice, minAmount, maxDiscount, posDiscount, minBonusAmonut float64, coupons []string) {

	//
	cheque.DeleteLoyaltyDiscounts() // remove all discounts
	cheque.DeleteMessages()         // remove all messages
	// TODO Добавить очистку оплат бонусами
	//cheque.DeleteBonusPayment() // remove all bonus payment

	// Setup new cheque line
	var line = models.Line{}
	line.ChequeLineNo = len(cheque.ChequeLines.ChequeLines) + 1
	line.Item.ItemUID = item
	line.Price = price
	line.Quantity = quantity
	line.Amount = price * quantity
	line.MinPrice = minPrice
	line.MinAmount = minAmount
	line.MinAmountAfterCurrencyDiscount = minBonusAmonut
	line.MaxDiscount = maxDiscount
	if len(coupons) > 0 {
		for _, coupon := range coupons {
			line.Coupon = append(line.Coupon, models.Coupon{CouponNo: coupon})
		}
	}
	// Setup pos discount
	if posDiscount != 0 {
		line.Discounts.Discounts = append(line.Discounts.Discounts, models.Discount{
			DiscountID: 0,
			Amount:     posDiscount,
		})
		line.Amount = line.Amount - posDiscount
	}

	// Setup new header
	cheque.PositionCount++
	cheque.Amount += line.Amount
	cheque.Amount = math.Round(cheque.Amount*100) / 100

	// Add new cheque line
	cheque.ChequeLines.ChequeLines = append(cheque.ChequeLines.ChequeLines, line)
}

func (cheque *models.Cheque) UpdateLine(number int, line models.Line) {
	for i, cl := range cheque.ChequeLines.ChequeLines {
		if cl.ChequeLineNo == number {
			// Clear loyalty
			cheque.DeleteLoyaltyDiscounts()
			cheque.DeleteMessages()

			// Update line
			cheque.ChequeLines.ChequeLines[i] = line

			//	Update receipt header
			cheque.Amount = 0
			for _, l := range cheque.ChequeLines.ChequeLines {
				cheque.Amount += l.Amount
			}

			break
		}
	}
}

func (cheque *models.Cheque) ChangeLine(number int, item string, price, quantity, minPrice, minAmount, maxDiscount, posDiscount, minBonusAmount float64, coupons []string) {
	//
	cheque.DeleteLoyaltyDiscounts() // remove all discounts
	cheque.DeleteMessages()         // remove all messages
	// TODO Добавить очистку оплат бонусами
	//cheque.DeleteBonusPayment() // remove all bonus payment

	//
	cheque.ChequeLines.ChequeLines[number].Item.ItemUID = item
	cheque.ChequeLines.ChequeLines[number].Price = price
	cheque.ChequeLines.ChequeLines[number].Quantity = quantity
	cheque.ChequeLines.ChequeLines[number].Amount = price * quantity
	cheque.ChequeLines.ChequeLines[number].MinPrice = math.Round(minPrice*100) / 100
	cheque.ChequeLines.ChequeLines[number].MinAmount = minAmount
	cheque.ChequeLines.ChequeLines[number].MinAmountAfterCurrencyDiscount = minBonusAmount
	cheque.ChequeLines.ChequeLines[number].MaxDiscount = maxDiscount
	cheque.ChequeLines.ChequeLines[number].Coupon = []models.Coupon{}
	for _, coupon := range coupons {
		cheque.ChequeLines.ChequeLines[number].Coupon = append(cheque.ChequeLines.ChequeLines[number].Coupon, models.Coupon{CouponNo: coupon})
	}

	//
	cheque.ChequeLines.ChequeLines[number].Discounts.Discounts = []models.Discount{}
	if posDiscount != 0 {
		cheque.ChequeLines.ChequeLines[number].Discounts.Discounts = append(cheque.ChequeLines.ChequeLines[number].Discounts.Discounts, models.Discount{
			DiscountID: 0,
			Amount:     posDiscount,
		})
		cheque.ChequeLines.ChequeLines[number].Amount = cheque.ChequeLines.ChequeLines[number].Amount - posDiscount
	}

	cheque.Amount = 0
	for _, line := range cheque.ChequeLines.ChequeLines {
		cheque.Amount += line.Amount
	}
	cheque.Amount = math.Round(cheque.Amount*100) / 100
}

func (cheque *models.Cheque) DeleteLine(index int) {
	//log.Trace().Int("lines", len(cheque.ChequeLines.ChequeLines)).Float64("amount", cheque.Amount).Msg("Performing deleting line from receipt:")

	//log.Trace().Msg("Cleaning loyalty discount from receipt")

	var count = 0
	var newChequeLine []models.Line

	// Adding remaining lines to new instance
	for _, line := range cheque.ChequeLines.ChequeLines {
		if line.ChequeLineNo != index {
			count += 1
			line.ChequeLineNo = count
			newChequeLine = append(newChequeLine, line)
		}
	}
	// Set new
	cheque.ChequeLines.ChequeLines = newChequeLine
	//
	cheque.Amount = 0
	for _, line := range cheque.ChequeLines.ChequeLines {
		cheque.Amount += line.Amount
	}
	cheque.PositionCount = int32(len(cheque.ChequeLines.ChequeLines))
}

//func (line *Line) AddPosDiscount(amount float64) {
//	td := cheque_discount.Amount
//	line.Discounts.Discounts = append(line.Discounts.Discounts, Discount{
//		DiscountID: 0,
//		Type:       &td,
//		Amount:     amount,
//	})
//}
//
//func (line *Line) AddCoupon(coupon string) {
//	line.Coupon = append(line.Coupon, Coupon{CouponNo: coupon})
//}
//
//func (line *Line) GetCoupons() (coupons string) {
//	for _, coupon := range line.Coupon {
//		coupons += fmt.Sprintf("%s,", coupon.CouponNo)
//	}
//	if coupons != "" {
//		coupons = coupons[:len(coupons)-1]
//	}
//	return coupons
//}

// GetLoyaltyPosDiscount Returns the loyalty discount in the specified receipt line
//func (line *Line) GetLoyaltyPosDiscount() (discounts float64) {
//	for _, discount := range line.Discounts.Discounts {
//		if discount.DiscountID != 0 && discount.Amount != 0 {
//			discounts += discount.Amount
//		}
//	}
//	return discounts
//}
//
//// GetPosDiscount Returns the cashier's discount in the specified receipt line
//func (line *Line) GetPosDiscount() (discounts float64) {
//	for _, discount := range line.Discounts.Discounts {
//		if discount.DiscountID == 0 && discount.Amount != 0 {
//			discounts += discount.Amount
//		}
//	}
//	return discounts
//}
