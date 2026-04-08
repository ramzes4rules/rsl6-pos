package models

// CardBalanceResult represents the result of a GetCardBalance call.
type CardBalanceResult struct {
	Balance struct {
		Value float32 `xml:"Value,attr"`
	}
	Msg struct {
		Device int    `xml:"Device,attr"`
		Body   string `xml:"Body,attr"`
	}
}
