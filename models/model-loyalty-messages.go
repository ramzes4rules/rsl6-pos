package models

import "encoding/xml"

type LoyaltyMessages struct {
	XMLName  xml.Name `xml:"Messages"`
	Messages []struct {
		XMLName   xml.Name   `xml:"Msg"`
		MessageID string     `xml:"MessageId,attr"`
		Device    DeviceType `xml:"Device,attr"`
		Body      string     `xml:"Body,attr"`
	} `xml:"Msg"`
}
