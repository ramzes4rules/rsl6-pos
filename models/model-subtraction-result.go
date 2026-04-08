package models

import "encoding/xml"

type Subtractions []SubtractedChequeLine

type SubtractedChequeLine struct {
	XMLName      xml.Name `xml:"ChequeLine"`
	ChequeLineNo int      `xml:"ChequeLineNo,attr"`
	Amount       float64  `xml:"Amount,attr"`
}
