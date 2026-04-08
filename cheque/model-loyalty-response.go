package cheque

import "encoding/xml"

// DiscountType represents the type of a loyalty discount.
type DiscountType string

const (
	DiscountPercent  DiscountType = "Percent"
	DiscountAmount   DiscountType = "Amount"
	DiscountTender   DiscountType = "Tender"
	DiscountFixPrice DiscountType = "FixPrice"
)

// LoyaltyDiscounts is the XML response from the GetDiscounts loyalty service call.
type LoyaltyDiscounts struct {
	XMLName     xml.Name `xml:"LoyaltyDiscounts"`
	ChequeLines struct {
		XMLName    xml.Name `xml:"ChequeLines"`
		ChequeLine []struct {
			XMLName      xml.Name `xml:"ChequeLine"`
			ChequeLineNo int      `xml:"ChequeLineNo,attr"`
			TotalAmount  float32  `xml:"TotalAmount,attr"`
			Discounts    struct {
				XMLName  xml.Name `xml:"Discounts"`
				Discount []struct {
					XMLName    xml.Name     `xml:"Discount"`
					DiscountID int          `xml:"DiscountID,attr"`
					Type       DiscountType `xml:"Type,attr"`
					Percent    float64      `xml:"Percent,attr"`
					Amount     float64      `xml:"Amount,attr"`
				} `xml:"Discount"`
			} `xml:"Discounts"`
		} `xml:"ChequeLine"`
	} `xml:"ChequeLines"`
}

// LoyaltyChequeLine represents a single cheque line within a LoyaltyDiscounts response.
//type LoyaltyChequeLine struct {
//	XMLName      xml.Name          `xml:"ChequeLine"`
//	ChequeLineNo int               `xml:"ChequeLineNo,attr"`
//	TotalAmount  float32           `xml:"TotalAmount,attr"`
//	Discounts    struct {
//		XMLName  xml.Name          `xml:"Discounts"`
//		Discount []struct {
//			XMLName    xml.Name     `xml:"Discount"`
//			DiscountID int          `xml:"DiscountID,attr"`
//			Type       DiscountType `xml:"Type,attr"`
//			Percent    float64      `xml:"Percent,attr"`
//			Amount     float64      `xml:"Amount,attr"`
//		}  `xml:"Discount"`
//	} `xml:"Discounts"`
//}

// IncomeDiscount represents a single discount within a loyalty service response.
//type IncomeDiscount struct {
//	XMLName    xml.Name     `xml:"Discount"`
//	DiscountID int          `xml:"DiscountID,attr"`
//	Type       DiscountType `xml:"Type,attr"`
//	Percent    float64      `xml:"Percent,attr"`
//	Amount     float64      `xml:"Amount,attr"`
//}

// LoyaltyMessages is the XML response from the GetMessages loyalty service call.
type LoyaltyMessages struct {
	XMLName  xml.Name `xml:"Messages"`
	Messages []struct {
		XMLName   xml.Name   `xml:"Msg"`
		MessageID string     `xml:"MessageId,attr"`
		Device    DeviceType `xml:"Device,attr"`
		Body      string     `xml:"Body,attr"`
	} `xml:"Msg"`
}

// LoyaltyMessage represents a single message within a LoyaltyMessages response.
//

// SubtractedChequeLine represents a single line within a SubtractBonus45 response.
type SubtractedChequeLine struct {
	XMLName      xml.Name `xml:"ChequeLine"`
	ChequeLineNo int      `xml:"ChequeLineNo,attr"`
	Amount       float64  `xml:"Amount,attr"`
}

// Subtractions is a list of subtracted cheque lines returned by SubtractBonus45.
type Subtractions []SubtractedChequeLine

// CardBalanceResult is the XML response from the GetCardBalance loyalty service call.
type CardBalanceResult struct {
	Balance struct {
		Value float32 `xml:"Value,attr"`
	}
	Msg struct {
		Device int    `xml:"Device,attr"`
		Body   string `xml:"Body,attr"`
	}
}
