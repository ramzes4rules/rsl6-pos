package cheque

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Type string
type Status string
type DeviceType string

const (
	TypeSale        Type       = "Sale"     // Тип чека продажи
	TypeReturn      Type       = "Return"   // Тип чека возврата
	StatusClosed    Status     = "Closed"   //
	StatusOpen      Status     = "Open"     //
	StatusCancelled Status     = "Canceled" //
	Slip            DeviceType = "0"        //
	CustomerDisplay DeviceType = "1"        //
	CashierDisplay  DeviceType = "2"        //
)

var CurrentCheque = new(Cheque)

type Cheque struct {
	XmlnsXsd        string         `xml:"xmlns:xsd,attr"`         //:xsd="http://www.w3.org/2001/XMLSchema"
	XmlnsXsi        string         `xml:"xmlns:xsi,attr"`         //:xsi="http://www.w3.org/2001/XMLSchema-instance"
	StoreID         uint8          `xml:"StoreID,attr"`           // Идентификатор магазина
	ShiftNo         string         `xml:"ShiftNo,attr"`           // Номер смены. Не используется.
	ChequeUID       string         `xml:"ChequeUID,attr"`         // Обязательный уникальный идентификатор чека
	ChequeNo        string         `xml:"ChequeNo,attr"`          // Номер чека на кассе
	OpenTime        string         `xml:"OpenTime,attr"`          // Время открытия чека
	CloseTime       string         `xml:"CloseTime,attr"`         // Время закрытия чека
	Amount          float64        `xml:"Amount,attr"`            // Сумма чека
	SubtractedBonus float64        `xml:"SubtractedBonus,attr"`   // Текущий вычет бонусов из чека в денежных единицах
	PositionCount   int32          `xml:"PositionCount,attr"`     // Количество позиций в чеке
	Status          Status         `xml:"Status,attr"`            // Статус чека
	ChequeType      Type           `xml:"ChequeType,attr"`        // Тип чека
	SaleCheque      *SaleCheque    `xml:"SaleCheque,omitempty"`   // Чек продажи
	DiscountCard    []DiscountCard `xml:"DiscountCard,omitempty"` // Описание данных по дисконтным картам
	Coupon          []Coupon       `xml:"Coupon,omitempty"`       // Описание данных купонов
	ChequeLines     Lines          `xml:"ChequeLines"`            // List of chequelines
	Discounts       *Discounts     `xml:"Discounts,omitempty"`    // Описание данных по скидкам для этой позиции
	Messages        *Messages      `xml:"Messages,omitempty"`     // Список сообщений в чеке
	Payments        *Payments      `xml:"Payment,omitempty"`      // List of payment
}

type SaleCheque struct {
	ChequeUID string `xml:"ChequeUID,attr,omitempty"`
	ChequeNo  string `xml:"ChequeNo,attr,omitempty"`
	OpenTime  string `xml:"OpenTime,attr,omitempty"`
}

type DiscountCard struct {
	DiscountCardNo       string  `xml:"DiscountCardNo,attr"`
	SubtractAmount       float64 `xml:"SubtractAmount,attr"`
	BonusCard            bool    `xml:"BonusCard,attr"`
	EnteredAsPhoneNumber bool    `xml:"EnteredAsPhoneNumber,attr"`
	SubtractedBonus      float64 `xml:"SubtractedBonus,attr,omitempty"`
}

type Coupon struct {
	CouponNo string `xml:"CouponNo,attr"` // Номер купона. Должен быть заведен в системе RS.Loyalty.
}

type Lines struct {
	XMLName     xml.Name `xml:"ChequeLines"`
	ChequeLines []Line   `xml:"ChequeLine"`
}

type Line struct {
	XMLName                        xml.Name  `xml:"ChequeLine"`
	ChequeLineNo                   int       `xml:"ChequeLineNo,attr"`                   //
	Price                          float64   `xml:"Price,attr"`                          // Цена товара по позиции
	Quantity                       float64   `xml:"Quantity,attr"`                       // Количество товара по позиции
	Amount                         float64   `xml:"Amount,attr"`                         // Итоговая сумма позиции БЕЗ учета скидки бонусами, но с учетом процентных скидок.
	MinAmount                      float64   `xml:"MinAmount,attr"`                      // Минимальная сумма по позиции
	MinPrice                       float64   `xml:"MinPrice,attr"`                       // Минимальная цена для позиции
	MaxDiscount                    float64   `xml:"MaxDiscount,attr"`                    // Максимальная скидка по позиции в процентах
	BonusDiscount                  float64   `xml:"BonusDiscount,attr"`                  // Сумма скидки бонусами по позиции примененная после списания бонусов
	MinAmountAfterCurrencyDiscount float64   `xml:"MinAmountAfterCurrencyDiscount,attr"` // Minimal receipt amount after bonus payment
	Item                           Item      `xml:"Item"`                                // Товар
	Discounts                      Discounts `xml:"Discounts,omitempty"`                 // Описание данных по скидкам для этой позиции
	Coupon                         []Coupon  `xml:"Coupon,omitempty"`                    // Примененный купон
}

type Discounts struct {
	Discounts []Discount `xml:"Discount,omitempty"`
}

type Discount struct {
	DiscountID int32         `xml:"DiscountID,attr"` // Внутренний ID скидки.
	Type       *DiscountType `xml:"Type,attr"`       // Тип скидки:
	Percent    *float64      `xml:"Percent,attr"`    // Величина процента скидки
	Amount     float64       `xml:"Amount,attr"`     // Величина скидки в единицах оплаты
}

type Item struct {
	ItemID  *int64  `xml:"ItemID,attr,omitempty"`  // Внутренний идентификатор товара (RS.Loyalty). В случае если этот код неизвестен, необходимо устанавливать значение 0.
	ItemUID string  `xml:"ItemUID,attr"`           // Уникальный идентификатор товара во внешней системе. 	Формат может быть любым. Максимум 50 символов.
	Barcode *string `xml:"Barcode,attr,omitempty"` // Штрих код товара.
}

type Messages struct {
	Messages []Message `xml:"Text,omitempty"`
}

type Message struct {
	XMLName   xml.Name   `xml:"Message"`
	MessageID string     `xml:"MessageID,attr"` // Внутренний ID сообщения.
	Device    DeviceType `xml:"Device,attr"`    // Тип скидки:
	Body      string     `xml:"Body,attr"`      // Текст сообщения.
}

type Payments struct {
	Payments []struct{} `xml:"Payment,omitempty"`
}

// AddPosDiscount adds a POS discount to the cheque line.
func (line *Line) AddPosDiscount(amount float64) {
	dt := DiscountAmount
	line.Discounts.Discounts = append(line.Discounts.Discounts, Discount{
		DiscountID: 0,
		Type:       &dt,
		Amount:     amount,
	})
}

// AddCoupon appends a coupon to the cheque line.
func (line *Line) AddCoupon(coupon string) {
	line.Coupon = append(line.Coupon, Coupon{CouponNo: coupon})
}

// GetPosDiscount returns the total POS discount applied to the cheque line.
func (line *Line) GetPosDiscount() float64 {
	var total float64
	for _, d := range line.Discounts.Discounts {
		if d.DiscountID == 0 {
			total += d.Amount
		}
	}
	return total
}

func newChequeNo() string {
	const set1 = "0123456789"

	var number = ""
	for i := 1; i <= 5; i++ {
		number += string(set1[rand.Intn(len(set1)-1)])
	}
	return number
}

func (cheque *Cheque) setChequeHeaders() {
	// Setup init values
	cheque.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	cheque.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	cheque.ChequeNo = newChequeNo()
	cheque.ChequeUID = uuid.NewString()
	cheque.Status = StatusOpen
	cheque.OpenTime = time.Now().Format(time.RFC3339)
	cheque.CloseTime = time.Now().Format(time.RFC3339)
	cheque.ChequeType = TypeSale
	cheque.ShiftNo = newChequeNo()
	cheque.StoreID = 1
	cheque.SaleCheque = nil
}

// CreateSaleReceipt setup init values and clears messages, discount, bonus payment. Don't clear check lines
func (cheque *Cheque) CreateSaleReceipt() {
	*cheque = Cheque{}
	cheque.setChequeHeaders()
}

// CopySaleReceipt df
func (cheque *Cheque) CopySaleReceipt() {
	cheque.setChequeHeaders()
	cheque.DeleteLoyaltyDiscounts()
	cheque.DeleteMessages()
	cheque.DeleteBonusPayment()
}

// CreateReturnReceipt se
func (cheque *Cheque) CreateReturnReceipt(number string, date string) {

	// Вызываем конструктор
	cheque.CreateSaleReceipt()

	//
	cheque.ChequeType = TypeReturn

	//
	cheque.SaleCheque = &SaleCheque{}
	cheque.SaleCheque.ChequeNo = number
	cheque.SaleCheque.OpenTime = date
}

// ModifyToReturnReceipt модифицирует чек продажи в чек возврата
func (cheque *Cheque) ModifyToReturnReceipt(variant int) {

	// Модифицируем заголовок чека
	cheque.setChequeHeaders()

	// Добавляем связь с чеком продажи
	cheque.SaleCheque = &SaleCheque{ChequeUID: cheque.ChequeUID}

	// Меняем строки чека
	for i, line := range cheque.ChequeLines.ChequeLines {
		cheque.ChequeLines.ChequeLines[i].Price = (line.Price*line.Quantity - cheque.GetLineLoyaltyDiscount(i) -
			cheque.GetLinePosDiscount(i) - line.BonusDiscount) / line.Quantity // Новая цена товара
		cheque.ChequeLines.ChequeLines[i].Amount = cheque.ChequeLines.ChequeLines[i].Price * line.Quantity // Новая сумма по позиции
		cheque.ChequeLines.ChequeLines[i].BonusDiscount = 0                                                // Сбрасываем сумму оплаты бонусами
		cheque.ChequeLines.ChequeLines[i].Discounts = Discounts{}                                          // Сбрасываем все скидки
	}

}

// SetOpenTime устанавливает время открытия чека
func (cheque *Cheque) SetOpenTime(t time.Time) {
	dt, err := time.Parse(time.RFC3339, cheque.OpenTime)
	if err != nil {
		dt = time.Now()
	}
	hour, minute, second := t.Clock()
	nd := time.Date(dt.Year(), dt.Month(), dt.Day(), hour, minute, second, dt.Nanosecond(), dt.Location())
	cheque.OpenTime = nd.Format(time.RFC3339)
}

// GetOpenTime возращает время открытия чека
func (cheque *Cheque) GetOpenTime() *time.Time {
	if ot, err := time.Parse(time.RFC3339, cheque.OpenTime); err == nil {
		return &ot
	}
	return nil
}

func (cheque *Cheque) SetOpenDate(t time.Time) {
	var dt, err = time.Parse(time.RFC3339, cheque.OpenTime)
	if err != nil {
		dt = time.Now()
	}
	year, month, day := t.Date()
	nd := time.Date(year, month, day, dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())
	cheque.OpenTime = fmt.Sprintf("%s", nd.Format(time.RFC3339))
}

func (cheque *Cheque) SetCloseTime(t time.Time) {
	var dt, err = time.Parse(time.RFC3339, cheque.OpenTime)
	if err != nil {
		dt = time.Now()
	}
	hour, minute, second := t.Clock()
	nd := time.Date(dt.Year(), dt.Month(), dt.Day(), hour, minute, second, dt.Nanosecond(), dt.Location())
	cheque.CloseTime = fmt.Sprintf("%s", nd.Format(time.RFC3339))
}

func (cheque *Cheque) SetCloseDate(t time.Time) {
	var dt, err = time.Parse(time.RFC3339, cheque.OpenTime)
	if err != nil {
		dt = time.Now()
	}
	year, month, day := t.Date()
	nd := time.Date(year, month, day, dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())
	cheque.CloseTime = fmt.Sprintf("%s", nd.Format(time.RFC3339))
}

// ParseFromXML parses the XML representation of the receipt and
// populates the check struct with the corresponding data.
// It returns an error if the parsing fails.
func (cheque *Cheque) ParseFromXML(text string) error {
	if err := xml.Unmarshal([]byte(text), cheque); err != nil {
		return err
	}
	return nil
}

// SerializeToXml returns the XML representation of the receipt.
// If indent is true, the XML will be formatted with indentation for better readability.
func (cheque *Cheque) SerializeToXml(indent bool) string {
	var response []byte
	if indent {
		response, _ = xml.MarshalIndent(cheque, "", " ")
	} else {
		response, _ = xml.Marshal(cheque)
	}
	result := append([]byte(xml.Header), response...)
	return string(result)
}

// CloseReceipt ggd
func (cheque *Cheque) CloseReceipt() {
	cheque.CloseTime = fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
	cheque.Status = StatusClosed
}

// CancelReceipt sets the receipt status to closed and updates the close time to the current time
func (cheque *Cheque) CancelReceipt() {
	cheque.CloseTime = fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
	cheque.Status = StatusCancelled
}

// GetSample populates the check with sample data for testing purposes.
func (cheque *Cheque) GetSample() *Cheque {

	cheque.setChequeHeaders()

	cheque.DiscountCard = append(cheque.DiscountCard, DiscountCard{
		DiscountCardNo:       "3846656766",
		SubtractAmount:       0,
		BonusCard:            true,
		EnteredAsPhoneNumber: false,
	})

	cheque.Coupon = append(cheque.Coupon, Coupon{CouponNo: "9900003221"})

	cheque.Messages = &Messages{
		Messages: []Message{{
			MessageID: "556",
			Device:    Slip,
			Body:      "Сообщение на слипе",
		}, {
			MessageID: "567",
			Device:    CashierDisplay,
			Body:      "Сообщение на экране кассира",
		}},
	}

	return cheque
}
