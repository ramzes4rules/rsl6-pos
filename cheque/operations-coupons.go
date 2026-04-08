package cheque

// AddCoupons sets the coupons on the receipt from the provided list.
func (cheque *Cheque) AddCoupons(coupons []string) {
	if len(coupons) > 0 {
		cheque.Coupon = []Coupon{}
		for _, coupon := range coupons {
			if coupon != "" {
				cheque.Coupon = append(cheque.Coupon, Coupon{CouponNo: coupon})
			}
		}
	} else {
		cheque.Coupon = []Coupon{}
	}
}

// GetReceiptCoupons returns the list of coupons applied to the receipt.
func (cheque *Cheque) GetReceiptCoupons() []Coupon {
	return cheque.Coupon
}
