package client

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"github.com/ramzes4rules/rsl6.pos.client/internal/soap"
)

// ==================== Ping ====================

func (c *Client) Ping(ctx context.Context) error {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ Ping"`
	}
	return c.call(ctx, "Ping", req{}, nil)
}

// ==================== IsOnline ====================

func (c *Client) IsOnline(ctx context.Context, version string) (bool, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ IsOnline"`
		Version string   `xml:"http://tempuri.org/ version"`
	}
	type resp struct {
		XMLName xml.Name `xml:"IsOnlineResponse"`
		Result  bool     `xml:"IsOnlineResult"`
	}
	var r resp
	if err := c.call(ctx, "IsOnline", req{Version: version}, &r); err != nil {
		return false, err
	}
	return r.Result, nil
}

// ==================== GetParameters ====================

func (c *Client) GetParameters(ctx context.Context, isStore bool) (map[string]string, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ GetParameters"`
		IsStore bool     `xml:"http://tempuri.org/ isStore"`
	}
	type kv struct {
		Key   string `xml:"Key"`
		Value string `xml:"Value"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetParametersResponse"`
		Result  struct {
			Items []kv `xml:"KeyValueOfstringstring"`
		} `xml:"GetParametersResult"`
	}
	var r resp
	if err := c.call(ctx, "GetParameters", req{IsStore: isStore}, &r); err != nil {
		return nil, err
	}
	result := make(map[string]string, len(r.Result.Items))
	for _, item := range r.Result.Items {
		result[item.Key] = item.Value
	}
	return result, nil
}

// ==================== GetStoreSettings ====================

func (c *Client) GetStoreSettings(ctx context.Context, currentSettings *StoreConfig) (*StoreConfig, error) {
	type req struct {
		XMLName         xml.Name               `xml:"http://tempuri.org/ GetStoreSettings"`
		CurrentSettings *soap.XmlStoreConfigNS `xml:"http://tempuri.org/ currentSettings,omitempty"`
	}
	type resp struct {
		XMLName xml.Name             `xml:"GetStoreSettingsResponse"`
		Result  *soap.XmlStoreConfig `xml:"GetStoreSettingsResult"`
	}
	var r resp
	if err := c.call(ctx, "GetStoreSettings", req{CurrentSettings: toXMLStoreConfig(currentSettings)}, &r); err != nil {
		return nil, err
	}
	return fromXMLStoreConfig(r.Result), nil
}

// ==================== RegisterDiscountCard ====================

func (c *Client) RegisterDiscountCard(ctx context.Context, discountCardNumber string) error {
	type req struct {
		XMLName            xml.Name `xml:"http://tempuri.org/ RegisterDiscountCard"`
		DiscountCardNumber string   `xml:"http://tempuri.org/ discountCardNumber"`
	}
	return c.call(ctx, "RegisterDiscountCard", req{DiscountCardNumber: discountCardNumber}, nil)
}

// ==================== IsCardValid ====================

func (c *Client) IsCardValid(ctx context.Context, discountCardNumber string) (bool, error) {
	type req struct {
		XMLName            xml.Name `xml:"http://tempuri.org/ IsCardValid"`
		DiscountCardNumber string   `xml:"http://tempuri.org/ discountCardNumber"`
	}
	type resp struct {
		XMLName xml.Name `xml:"IsCardValidResponse"`
		Result  bool     `xml:"IsCardValidResult"`
	}
	var r resp
	if err := c.call(ctx, "IsCardValid", req{DiscountCardNumber: discountCardNumber}, &r); err != nil {
		return false, err
	}
	return r.Result, nil
}

// ==================== IsCouponValid ====================

func (c *Client) IsCouponValid(ctx context.Context, couponNumber string) (bool, error) {
	type req struct {
		XMLName      xml.Name `xml:"http://tempuri.org/ IsCouponValid"`
		CouponNumber string   `xml:"http://tempuri.org/ couponNumber"`
	}
	type resp struct {
		XMLName xml.Name `xml:"IsCouponValidResponse"`
		Result  bool     `xml:"IsCouponValidResult"`
	}
	var r resp
	if err := c.call(ctx, "IsCouponValid", req{CouponNumber: couponNumber}, &r); err != nil {
		return false, err
	}
	return r.Result, nil
}

// ==================== GetVerifyCode ====================

func (c *Client) GetVerifyCode(ctx context.Context, discountCardNumber string) (string, error) {
	type req struct {
		XMLName            xml.Name `xml:"http://tempuri.org/ GetVerifyCode"`
		DiscountCardNumber string   `xml:"http://tempuri.org/ discountCardNumber"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetVerifyCodeResponse"`
		Result  string   `xml:"GetVerifyCodeResult"`
	}
	var r resp
	if err := c.call(ctx, "GetVerifyCode", req{DiscountCardNumber: discountCardNumber}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== GetCardBalance ====================

func (c *Client) GetCardBalance(ctx context.Context, discountCardNumber string) (string, error) {
	type req struct {
		XMLName            xml.Name `xml:"http://tempuri.org/ GetCardBalance"`
		DiscountCardNumber string   `xml:"http://tempuri.org/ discountCardNumber"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetCardBalanceResponse"`
		Result  string   `xml:"GetCardBalanceResult"`
	}
	var r resp
	if err := c.call(ctx, "GetCardBalance", req{DiscountCardNumber: discountCardNumber}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== GetCardDiscountAmount ====================

func (c *Client) GetCardDiscountAmount(ctx context.Context, discountCardNumber, cheque string) (float64, error) {
	type req struct {
		XMLName            xml.Name `xml:"http://tempuri.org/ GetCardDiscountAmount"`
		DiscountCardNumber string   `xml:"http://tempuri.org/ discountCardNumber"`
		Cheque             string   `xml:"http://tempuri.org/ cheque"`
	}
	type resp struct {
		XMLName xml.Name        `xml:"GetCardDiscountAmountResponse"`
		Result  soap.XmlDecimal `xml:"GetCardDiscountAmountResult"`
	}
	var r resp
	if err := c.call(ctx, "GetCardDiscountAmount", req{DiscountCardNumber: discountCardNumber, Cheque: cheque}, &r); err != nil {
		return 0, err
	}
	return float64(r.Result), nil
}

// ==================== GetCardDiscountAmountString ====================

func (c *Client) GetCardDiscountAmountString(ctx context.Context, discountCardNumber, cheque string) (string, error) {
	type req struct {
		XMLName            xml.Name `xml:"http://tempuri.org/ GetCardDiscountAmountString"`
		DiscountCardNumber string   `xml:"http://tempuri.org/ discountCardNumber"`
		Cheque             string   `xml:"http://tempuri.org/ cheque"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetCardDiscountAmountStringResponse"`
		Result  string   `xml:"GetCardDiscountAmountStringResult"`
	}
	var r resp
	if err := c.call(ctx, "GetCardDiscountAmountString", req{DiscountCardNumber: discountCardNumber, Cheque: cheque}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== GetMessages ====================

func (c *Client) GetMessages(ctx context.Context, cheque string) (string, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ GetMessages"`
		Cheque  string   `xml:"http://tempuri.org/ cheque"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetMessagesResponse"`
		Result  string   `xml:"GetMessagesResult"`
	}
	var r resp
	if err := c.call(ctx, "GetMessages", req{Cheque: cheque}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== GetDiscounts ====================

func (c *Client) GetDiscounts(ctx context.Context, cheque string) (string, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ GetDiscounts"`
		Cheque  string   `xml:"http://tempuri.org/ cheque"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetDiscountsResponse"`
		Result  string   `xml:"GetDiscountsResult"`
	}
	var r resp
	if err := c.call(ctx, "GetDiscounts", req{Cheque: cheque}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== GetEmail ====================

func (c *Client) GetEmail(ctx context.Context, discountNumber string) (string, error) {
	type req struct {
		XMLName        xml.Name `xml:"http://tempuri.org/ GetEmail"`
		DiscountNumber string   `xml:"http://tempuri.org/ discountNumber"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetEmailResponse"`
		Result  string   `xml:"GetEmailResult"`
	}
	var r resp
	if err := c.call(ctx, "GetEmail", req{DiscountNumber: discountNumber}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== GetSelfBuyDiscounts ====================

func (c *Client) GetSelfBuyDiscounts(ctx context.Context, cheque string, storeID int64) (string, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ GetSelfBuyDiscounts"`
		Cheque  string   `xml:"http://tempuri.org/ cheque"`
		StoreID int64    `xml:"http://tempuri.org/ storeId"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetSelfBuyDiscountsResponse"`
		Result  string   `xml:"GetSelfBuyDiscountsResult"`
	}
	var r resp
	if err := c.call(ctx, "GetSelfBuyDiscounts", req{Cheque: cheque, StoreID: storeID}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== Accrual ====================

func (c *Client) Accrual(ctx context.Context, cheque string) (string, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ Accrual"`
		Cheque  string   `xml:"http://tempuri.org/ cheque"`
	}
	type resp struct {
		XMLName xml.Name `xml:"AccrualResponse"`
		Result  string   `xml:"AccrualResult"`
	}
	var r resp
	if err := c.call(ctx, "Accrual", req{Cheque: cheque}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== OfflineAccrual ====================

func (c *Client) OfflineAccrual(ctx context.Context, cheque string) (bool, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ OfflineAccrual"`
		Cheque  string   `xml:"http://tempuri.org/ cheque"`
	}
	type resp struct {
		XMLName xml.Name `xml:"OfflineAccrualResponse"`
		Result  bool     `xml:"OfflineAccrualResult"`
	}
	var r resp
	if err := c.call(ctx, "OfflineAccrual", req{Cheque: cheque}, &r); err != nil {
		return false, err
	}
	return r.Result, nil
}

// ==================== Refund ====================

func (c *Client) Refund(ctx context.Context, refundCheque string, chequeID int64) error {
	type req struct {
		XMLName      xml.Name `xml:"http://tempuri.org/ Refund"`
		RefundCheque string   `xml:"http://tempuri.org/ refundRefundCheque"`
		ChequeID     int64    `xml:"http://tempuri.org/ chequeID"`
	}
	return c.call(ctx, "Refund", req{RefundCheque: refundCheque, ChequeID: chequeID}, nil)
}

// ==================== SubtractBonus ====================

func (c *Client) SubtractBonus(ctx context.Context, discountCardNumber string, amount float64, cheque string) error {
	type req struct {
		XMLName            xml.Name        `xml:"http://tempuri.org/ SubtractBonus"`
		DiscountCardNumber string          `xml:"http://tempuri.org/ discountCardNumber"`
		Amount             soap.XmlDecimal `xml:"http://tempuri.org/ amount"`
		Cheque             string          `xml:"http://tempuri.org/ cheque"`
	}
	return c.call(ctx, "SubtractBonus", req{
		DiscountCardNumber: discountCardNumber,
		Amount:             soap.XmlDecimal(amount),
		Cheque:             cheque,
	}, nil)
}

// ==================== SubtractBonus45 ====================

func (c *Client) SubtractBonus45(ctx context.Context, discountCardNumber string, amount float64, cheque string) (string, error) {
	type req struct {
		XMLName            xml.Name        `xml:"http://tempuri.org/ SubtractBonus45"`
		DiscountCardNumber string          `xml:"http://tempuri.org/ discountCardNumber"`
		Amount             soap.XmlDecimal `xml:"http://tempuri.org/ amount"`
		Cheque             string          `xml:"http://tempuri.org/ cheque"`
	}
	type resp struct {
		XMLName xml.Name `xml:"SubtractBonus45Response"`
		Result  string   `xml:"SubtractBonus45Result"`
	}
	var r resp
	if err := c.call(ctx, "SubtractBonus45", req{
		DiscountCardNumber: discountCardNumber,
		Amount:             soap.XmlDecimal(amount),
		Cheque:             cheque,
	}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== CancelSubtractBonus ====================

func (c *Client) CancelSubtractBonus(ctx context.Context, discountCardNumber string, amount float64, cheque string) error {
	type req struct {
		XMLName            xml.Name        `xml:"http://tempuri.org/ CancelSubtractBonus"`
		DiscountCardNumber string          `xml:"http://tempuri.org/ discountCardNumber"`
		Amount             soap.XmlDecimal `xml:"http://tempuri.org/ amount"`
		Cheque             string          `xml:"http://tempuri.org/ cheque"`
	}
	return c.call(ctx, "CancelSubtractBonus", req{
		DiscountCardNumber: discountCardNumber,
		Amount:             soap.XmlDecimal(amount),
		Cheque:             cheque,
	}, nil)
}

// ==================== ValidateUser ====================

func (c *Client) ValidateUser(ctx context.Context, userName, password string) (bool, error) {
	type req struct {
		XMLName  xml.Name `xml:"http://tempuri.org/ ValidateUser"`
		UserName string   `xml:"http://tempuri.org/ userName"`
		Password string   `xml:"http://tempuri.org/ password"`
	}
	type resp struct {
		XMLName xml.Name `xml:"ValidateUserResponse"`
		Result  bool     `xml:"ValidateUserResult"`
	}
	var r resp
	if err := c.call(ctx, "ValidateUser", req{UserName: userName, Password: password}, &r); err != nil {
		return false, err
	}
	return r.Result, nil
}

// ==================== CheckDiscountCard ====================

func (c *Client) CheckDiscountCard(ctx context.Context, discountCard string) (bool, error) {
	type req struct {
		XMLName      xml.Name `xml:"http://tempuri.org/ CheckDiscountCard"`
		DiscountCard string   `xml:"http://tempuri.org/ discountCard"`
	}
	type resp struct {
		XMLName xml.Name `xml:"CheckDiscountCardResponse"`
		Result  bool     `xml:"CheckDiscountCardResult"`
	}
	var r resp
	if err := c.call(ctx, "CheckDiscountCard", req{DiscountCard: discountCard}, &r); err != nil {
		return false, err
	}
	return r.Result, nil
}

// ==================== ValidateUserRole ====================

func (c *Client) ValidateUserRole(ctx context.Context, userName, roleName string) (bool, error) {
	type req struct {
		XMLName  xml.Name `xml:"http://tempuri.org/ ValidateUserRole"`
		UserName string   `xml:"http://tempuri.org/ userName"`
		RoleName string   `xml:"http://tempuri.org/ roleName"`
	}
	type resp struct {
		XMLName xml.Name `xml:"ValidateUserRoleResponse"`
		Result  bool     `xml:"ValidateUserRoleResult"`
	}
	var r resp
	if err := c.call(ctx, "ValidateUserRole", req{UserName: userName, RoleName: roleName}, &r); err != nil {
		return false, err
	}
	return r.Result, nil
}

// ==================== GetUserRole ====================

func (c *Client) GetUserRole(ctx context.Context, userName string) (string, error) {
	type req struct {
		XMLName  xml.Name `xml:"http://tempuri.org/ GetUserRole"`
		UserName string   `xml:"http://tempuri.org/ userName"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetUserRoleResponse"`
		Result  string   `xml:"GetUserRoleResult"`
	}
	var r resp
	if err := c.call(ctx, "GetUserRole", req{UserName: userName}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== ActivationPaymentCard ====================

func (c *Client) ActivationPaymentCard(ctx context.Context, discountCard string) (float64, error) {
	type req struct {
		XMLName      xml.Name `xml:"http://tempuri.org/ ActivationPaymentCard"`
		DiscountCard string   `xml:"http://tempuri.org/ discountCard"`
	}
	type resp struct {
		XMLName xml.Name        `xml:"ActivationPaymentCardResponse"`
		Result  soap.XmlDecimal `xml:"ActivationPaymentCardResult"`
	}
	var r resp
	if err := c.call(ctx, "ActivationPaymentCard", req{DiscountCard: discountCard}, &r); err != nil {
		return 0, err
	}
	return float64(r.Result), nil
}

// ==================== CancelActivationPaymentCard ====================

func (c *Client) CancelActivationPaymentCard(ctx context.Context, discountCard string) (bool, error) {
	type req struct {
		XMLName      xml.Name `xml:"http://tempuri.org/ CancelActivationPaymentCard"`
		DiscountCard string   `xml:"http://tempuri.org/ discountCard"`
	}
	type resp struct {
		XMLName xml.Name `xml:"CancelActivationPaymentCardResponse"`
		Result  bool     `xml:"CancelActivationPaymentCardResult"`
	}
	var r resp
	if err := c.call(ctx, "CancelActivationPaymentCard", req{DiscountCard: discountCard}, &r); err != nil {
		return false, err
	}
	return r.Result, nil
}

// ==================== QuerySyncStream ====================

func (c *Client) QuerySyncStream(ctx context.Context, data []TupleOfStringLong) (string, error) {
	type tuple struct {
		XMLName xml.Name `xml:"http://schemas.datacontract.org/2004/07/System TupleOfstringlong"`
		Item1   string   `xml:"http://schemas.datacontract.org/2004/07/System m_Item1"`
		Item2   int64    `xml:"http://schemas.datacontract.org/2004/07/System m_Item2"`
	}
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ QuerySyncStream"`
		Data    struct {
			Items []tuple
		} `xml:"http://tempuri.org/ data"`
	}
	r := req{}
	for _, d := range data {
		r.Data.Items = append(r.Data.Items, tuple{Item1: d.Item1, Item2: d.Item2})
	}
	type resp struct {
		XMLName xml.Name `xml:"QuerySyncStreamResponse"`
		Result  string   `xml:"QuerySyncStreamResult"`
	}
	var result resp
	if err := c.call(ctx, "QuerySyncStream", r, &result); err != nil {
		return "", err
	}
	return result.Result, nil
}

// ==================== GetSyncStream ====================

func (c *Client) GetSyncStream(ctx context.Context, taskID string) ([]byte, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ GetSyncStream"`
		TaskID  string   `xml:"http://tempuri.org/ taskID"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetSyncStreamResponse"`
		Result  string   `xml:"GetSyncStreamResult"`
	}
	var r resp
	if err := c.call(ctx, "GetSyncStream", req{TaskID: taskID}, &r); err != nil {
		return nil, err
	}
	if r.Result == "" {
		return nil, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(r.Result)
	if err != nil {
		return nil, fmt.Errorf("decode sync stream base64: %w", err)
	}
	return decoded, nil
}

// ==================== IsTaskCompleted ====================

func (c *Client) IsTaskCompleted(ctx context.Context, taskID string) (bool, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ IsTaskCompleted"`
		TaskID  string   `xml:"http://tempuri.org/ taskID"`
	}
	type resp struct {
		XMLName xml.Name `xml:"IsTaskCompletedResponse"`
		Result  bool     `xml:"IsTaskCompletedResult"`
	}
	var r resp
	if err := c.call(ctx, "IsTaskCompleted", req{TaskID: taskID}, &r); err != nil {
		return false, err
	}
	return r.Result, nil
}

// ==================== GetUpdateStream ====================

func (c *Client) GetUpdateStream(ctx context.Context, filename string) ([]byte, error) {
	type req struct {
		XMLName  xml.Name `xml:"http://tempuri.org/ GetUpdateStream"`
		Filename string   `xml:"http://tempuri.org/ filename"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetUpdateStreamResponse"`
		Result  string   `xml:"GetUpdateStreamResult"`
	}
	var r resp
	if err := c.call(ctx, "GetUpdateStream", req{Filename: filename}, &r); err != nil {
		return nil, err
	}
	if r.Result == "" {
		return nil, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(r.Result)
	if err != nil {
		return nil, fmt.Errorf("decode update stream base64: %w", err)
	}
	return decoded, nil
}

// ==================== UploadReferences ====================

func (c *Client) UploadReferences(ctx context.Context, packet string, stamp int64) error {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ UploadReferences"`
		Packet  string   `xml:"http://tempuri.org/ packet"`
		Stamp   int64    `xml:"http://tempuri.org/ stamp"`
	}
	return c.call(ctx, "UploadReferences", req{Packet: packet, Stamp: stamp}, nil)
}

// ==================== GetReferencesStamp ====================

func (c *Client) GetReferencesStamp(ctx context.Context) (int64, error) {
	type req struct {
		XMLName xml.Name `xml:"http://tempuri.org/ GetReferencesStamp"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetReferencesStampResponse"`
		Result  int64    `xml:"GetReferencesStampResult"`
	}
	var r resp
	if err := c.call(ctx, "GetReferencesStamp", req{}, &r); err != nil {
		return 0, err
	}
	return r.Result, nil
}

// ==================== GetDataPacket ====================

func (c *Client) GetDataPacket(ctx context.Context, paramPacket string) (string, error) {
	type req struct {
		XMLName     xml.Name `xml:"http://tempuri.org/ GetDataPacket"`
		ParamPacket string   `xml:"http://tempuri.org/ paramPacket"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetDataPacketResponse"`
		Result  string   `xml:"GetDataPacketResult"`
	}
	var r resp
	if err := c.call(ctx, "GetDataPacket", req{ParamPacket: paramPacket}, &r); err != nil {
		return "", err
	}
	return r.Result, nil
}

// ==================== SendInfoPacket ====================

func (c *Client) SendInfoPacket(ctx context.Context, packet *RSInfoPacket) error {
	type xmlPacket struct {
		OfflineChequeCount int64  `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol OfflineChequeCount"`
		Version            string `xml:"http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol Version"`
	}
	type req struct {
		XMLName xml.Name   `xml:"http://tempuri.org/ SendInfoPacket"`
		Packet  *xmlPacket `xml:"http://tempuri.org/ packet,omitempty"`
	}
	var p *xmlPacket
	if packet != nil {
		p = &xmlPacket{
			OfflineChequeCount: packet.OfflineChequeCount,
			Version:            packet.Version,
		}
	}
	return c.call(ctx, "SendInfoPacket", req{Packet: p}, nil)
}

// ==================== GetStatistic ====================

func (c *Client) GetStatistic(ctx context.Context, params GetStatisticRequest) ([]ItemStatistics, error) {
	type req struct {
		XMLName       xml.Name          `xml:"http://tempuri.org/ GetStatistic"`
		AccountID     int64             `xml:"http://tempuri.org/ accointId"`
		ItemIDs       soap.XmlLongArray `xml:"http://tempuri.org/ itemIds"`
		Time          soap.XmlDateTime  `xml:"http://tempuri.org/ time"`
		StatisticFlag string            `xml:"http://tempuri.org/ statisticFlag"`
	}
	ra := req{
		AccountID:     params.AccountID,
		ItemIDs:       soap.XmlLongArray{Items: params.ItemIDs},
		Time:          soap.XmlDateTime(params.Time),
		StatisticFlag: string(params.StatisticFlag),
	}
	type itemStat struct {
		DailyQuantity   soap.XmlDecimal `xml:"DailyQuantity"`
		ItemID          int64           `xml:"ItemId"`
		MonthlyQuantity soap.XmlDecimal `xml:"MonthlyQuantity"`
		WeeklyQuantity  soap.XmlDecimal `xml:"WeeklyQuantity"`
	}
	type resp struct {
		XMLName xml.Name `xml:"GetStatisticResponse"`
		Result  struct {
			Items []itemStat `xml:"ItemStatistics"`
		} `xml:"GetStatisticResult"`
	}
	var result resp
	if err := c.call(ctx, "GetStatistic", ra, &result); err != nil {
		return nil, err
	}
	stats := make([]ItemStatistics, 0, len(result.Result.Items))
	for _, item := range result.Result.Items {
		stats = append(stats, ItemStatistics{
			DailyQuantity:   float64(item.DailyQuantity),
			ItemID:          item.ItemID,
			MonthlyQuantity: float64(item.MonthlyQuantity),
			WeeklyQuantity:  float64(item.WeeklyQuantity),
		})
	}
	return stats, nil
}

// ==================== StoreConfig converters ====================

func toXMLStoreConfig(sc *StoreConfig) *soap.XmlStoreConfigNS {
	if sc == nil {
		return nil
	}
	return &soap.XmlStoreConfigNS{
		NoAddBonusForAdvertising:  sc.NoAddBonusForAdvertising,
		NoDiscountsForAdvertising: sc.NoDiscountsForAdvertising,
		NoPayBonusForAdvertising:  sc.NoPayBonusForAdvertising,
		OfflineCheckTime:          sc.OfflineCheckTime,
		OfflineChequeSendCount:    sc.OfflineChequeSendCount,
		OfflineDiscount:           sc.OfflineDiscount,
		OnlineCheckTime:           sc.OnlineCheckTime,
		StoreSettingsID:           sc.StoreSettingsID,
		SyncroTimeout:             sc.SyncroTimeout,
		Timeout:                   sc.Timeout,
		UseMapping:                sc.UseMapping,
	}
}

func fromXMLStoreConfig(x *soap.XmlStoreConfig) *StoreConfig {
	if x == nil {
		return nil
	}
	return &StoreConfig{
		NoAddBonusForAdvertising:  x.NoAddBonusForAdvertising,
		NoDiscountsForAdvertising: x.NoDiscountsForAdvertising,
		NoPayBonusForAdvertising:  x.NoPayBonusForAdvertising,
		OfflineCheckTime:          x.OfflineCheckTime,
		OfflineChequeSendCount:    x.OfflineChequeSendCount,
		OfflineDiscount:           x.OfflineDiscount,
		OnlineCheckTime:           x.OnlineCheckTime,
		StoreSettingsID:           x.StoreSettingsID,
		SyncroTimeout:             x.SyncroTimeout,
		Timeout:                   x.Timeout,
		UseMapping:                x.UseMapping,
	}
}
