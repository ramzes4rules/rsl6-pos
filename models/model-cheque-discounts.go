package models

import (
	"encoding/xml"
)

type DiscountType string

const (
	Percent  DiscountType = "Percent"
	Amount   DiscountType = "Amount"
	Tender   DiscountType = "Tender"
	FixPrice DiscountType = "FixPrice"
)

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

//type IncomeChequeLines struct {
//	XMLName    xml.Name           `xml:"ChequeLines"`
//	ChequeLine []struct {
//		XMLName      xml.Name        `xml:"ChequeLine"`
//		ChequeLineNo int             `xml:"ChequeLineNo,attr"`
//		TotalAmount  float32         `xml:"TotalAmount,attr"`
//		Discounts    struct {
//			XMLName  xml.Name         `xml:"Discounts"`
//			Discount []struct {
//				XMLName    xml.Name     `xml:"Discount"`
//				DiscountID int          `xml:"DiscountID,attr"`
//				Type       DiscountType `xml:"Type,attr"`
//				Percent    float64      `xml:"Percent,attr"`
//				Amount     float64      `xml:"Amount,attr"`
//			} `xml:"Discount"`
//		} `xml:"Discounts"`
//	} `xml:"ChequeLine"`
//}

//type IncomeChequeLine struct {
//	XMLName      xml.Name        `xml:"ChequeLine"`
//	ChequeLineNo int             `xml:"ChequeLineNo,attr"`
//	TotalAmount  float32         `xml:"TotalAmount,attr"`
//	Discounts    struct {
//		XMLName  xml.Name         `xml:"Discounts"`
//		Discount []struct {
//			XMLName    xml.Name     `xml:"Discount"`
//			DiscountID int          `xml:"DiscountID,attr"`
//			Type       DiscountType `xml:"Type,attr"`
//			Percent    float64      `xml:"Percent,attr"`
//			Amount     float64      `xml:"Amount,attr"`
//		} `xml:"Discount"`
//	} `xml:"Discounts"`
//}

//type IncomeDiscounts struct {
//	XMLName  xml.Name         `xml:"Discounts"`
//	Discount []struct {
//		XMLName    xml.Name     `xml:"Discount"`
//		DiscountID int          `xml:"DiscountID,attr"`
//		Type       DiscountType `xml:"Type,attr"`
//		Percent    float64      `xml:"Percent,attr"`
//		Amount     float64      `xml:"Amount,attr"`
//	} `xml:"Discount"`
//}

//type IncomeDiscount struct {
//	XMLName    xml.Name     `xml:"Discount"`
//	DiscountID int          `xml:"DiscountID,attr"`
//	Type       DiscountType `xml:"Type,attr"`
//	Percent    float64      `xml:"Percent,attr"`
//	Amount     float64      `xml:"Amount,attr"`
//}
