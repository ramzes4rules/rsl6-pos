package cheque

func (cheque *Cheque) AddCoupons(coupons []string) {
	//log.Trace().Msg("Adding coupons to receipt body")
	//cpns := strings.StringToArray(coupons)
	if len(coupons) > 0 {
		//log.Trace().Int("number", len(cpns)).Msg("Coupons found:")
		cheque.Coupon = []Coupon{}
		for _, coupon := range coupons {
			if coupon != "" {
				cheque.Coupon = append(cheque.Coupon, Coupon{CouponNo: coupon})
			}
		}
		//log.Trace().Int("coupons", len(cheque.Coupon)).Msg("All coupons were added:")
	} else {
		//log.Trace().Msg("Coupons not found")
		cheque.Coupon = []Coupon{}
	}
}

// GetReceiptCoupons returns the list of coupons applied to the receipt
func (cheque *Cheque) GetReceiptCoupons() []Coupon {
	return cheque.Coupon
}
