package soap

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ---------- XmlDecimal: xs:decimal without scientific notation ----------

// XmlDecimal handles xs:decimal serialization without scientific notation.
type XmlDecimal float64

func (d XmlDecimal) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s := strconv.FormatFloat(float64(d), 'f', -1, 64)
	return e.EncodeElement(s, start)
}

func (d *XmlDecimal) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := dec.DecodeElement(&s, &start); err != nil {
		return err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		*d = 0
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("parse decimal %q: %w", s, err)
	}
	*d = XmlDecimal(f)
	return nil
}

// ---------- XmlDateTime: xs:dateTime in WCF format ----------

// XmlDateTime handles xs:dateTime serialization in WCF format.
type XmlDateTime time.Time

const DateTimeFmt = "2006-01-02T15:04:05"

func (t XmlDateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s := time.Time(t).Format(DateTimeFmt)
	return e.EncodeElement(s, start)
}

func (t *XmlDateTime) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := dec.DecodeElement(&s, &start); err != nil {
		return err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parsed, err := time.Parse(DateTimeFmt, s)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return fmt.Errorf("parse datetime %q: %w", s, err)
		}
	}
	*t = XmlDateTime(parsed)
	return nil
}

// ---------- XmlStoreConfig: response parsing (lenient, no namespace) ----------

// XmlStoreConfig maps StoreConfig for XML deserialization (response).
type XmlStoreConfig struct {
	NoAddBonusForAdvertising  bool  `xml:"NoAddBonusForAdvertising"`
	NoDiscountsForAdvertising bool  `xml:"NoDiscountsForAdvertising"`
	NoPayBonusForAdvertising  bool  `xml:"NoPayBonusForAdvertising"`
	OfflineCheckTime          int64 `xml:"OfflineCheckTime"`
	OfflineChequeSendCount    int32 `xml:"OfflineChequeSendCount"`
	OfflineDiscount           bool  `xml:"OfflineDiscount"`
	OnlineCheckTime           int64 `xml:"OnlineCheckTime"`
	StoreSettingsID           int64 `xml:"StoreSettingsID"`
	SyncroTimeout             int64 `xml:"SyncroTimeout"`
	Timeout                   int64 `xml:"Timeout"`
	UseMapping                bool  `xml:"UseMapping"`
}

// ---------- XmlStoreConfigNS: request serialization (explicit namespace) ----------

// XmlStoreConfigNS maps StoreConfig for XML serialization with explicit namespace (requests).
type XmlStoreConfigNS struct {
	NoAddBonusForAdvertising  bool  `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol NoAddBonusForAdvertising"`
	NoDiscountsForAdvertising bool  `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol NoDiscountsForAdvertising"`
	NoPayBonusForAdvertising  bool  `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol NoPayBonusForAdvertising"`
	OfflineCheckTime          int64 `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol OfflineCheckTime"`
	OfflineChequeSendCount    int32 `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol OfflineChequeSendCount"`
	OfflineDiscount           bool  `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol OfflineDiscount"`
	OnlineCheckTime           int64 `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol OnlineCheckTime"`
	StoreSettingsID           int64 `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol StoreSettingsID"`
	SyncroTimeout             int64 `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol SyncroTimeout"`
	Timeout                   int64 `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol Timeout"`
	UseMapping                bool  `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol UseMapping"`
}

// ---------- XmlLongArray: ArrayOflong with proper namespace ----------

// XmlLongArray handles ArrayOflong serialization with proper namespace.
type XmlLongArray struct {
	Items []int64
}

func (a XmlLongArray) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for _, v := range a.Items {
		elem := xml.StartElement{
			Name: xml.Name{Space: NsArrays, Local: "long"},
		}
		if err := e.EncodeElement(v, elem); err != nil {
			return fmt.Errorf("encode long: %w", err)
		}
	}
	return e.EncodeToken(start.End())
}

