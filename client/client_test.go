package rslpos

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ---------- Test helpers ----------

// soapResp wraps body XML in a SOAP 1.2 response envelope.
func soapResp(body string) string {
	return `<?xml version="1.0" encoding="utf-8"?>` +
		`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">` +
		`<s:Body>` + body + `</s:Body></s:Envelope>`
}

// soapFaultResp creates a SOAP fault response.
func soapFaultResp(code, reason, detail string) string {
	d := ""
	if detail != "" {
		d = `<s:Detail>` + detail + `</s:Detail>`
	}
	return soapResp(`<s:Fault><s:Code><s:Value>` + code + `</s:Value></s:Code>` +
		`<s:Reason><s:Text>` + reason + `</s:Text></s:Reason>` + d + `</s:Fault>`)
}

// newTestClient creates an httptest.Server + Client. The handler validates SOAPAction
// and returns the given response.
func newTestClient(t *testing.T, expectedAction string, response string) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.Contains(ct, "application/soap+xml") {
			t.Errorf("expected soap+xml content type, got %s", ct)
		}
		if expectedAction != "" && !strings.Contains(ct, expectedAction) {
			t.Errorf("expected action %s in content type, got %s", expectedAction, ct)
		}
		// Read and discard request body
		_, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)
	return NewClient(server.URL)
}

// newTestClientWithValidator creates a test client that also validates the request body.
func newTestClientWithValidator(t *testing.T, expectedAction string, validator func(body string), response string) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if expectedAction != "" && !strings.Contains(ct, expectedAction) {
			t.Errorf("expected action %s in content type, got %s", expectedAction, ct)
		}
		body, _ := io.ReadAll(r.Body)
		if validator != nil {
			validator(string(body))
		}
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)
	return NewClient(server.URL)
}

func ctx() context.Context {
	return context.Background()
}

// ---------- Tests for all 36 operations ----------

func TestPing(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/Ping",
		soapResp(`<PingResponse xmlns="http://tempuri.org/"/>`))
	err := c.Ping(ctx())
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestIsOnline(t *testing.T) {
	c := newTestClientWithValidator(t, "IRSLoyaltyService/IsOnline",
		func(body string) {
			if !strings.Contains(body, "version") {
				t.Error("request should contain version element")
			}
		},
		soapResp(`<IsOnlineResponse xmlns="http://tempuri.org/"><IsOnlineResult>true</IsOnlineResult></IsOnlineResponse>`))
	result, err := c.IsOnline(ctx(), "2.0")
	if err != nil {
		t.Fatalf("IsOnline: %v", err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestIsOnline_False(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<IsOnlineResponse xmlns="http://tempuri.org/"><IsOnlineResult>false</IsOnlineResult></IsOnlineResponse>`))
	result, err := c.IsOnline(ctx(), "1.0")
	if err != nil {
		t.Fatalf("IsOnline: %v", err)
	}
	if result {
		t.Error("expected false")
	}
}

func TestGetParameters(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetParameters",
		soapResp(`<GetParametersResponse xmlns="http://tempuri.org/">
			<GetParametersResult xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
				<a:KeyValueOfstringstring><a:Key>param1</a:Key><a:Value>value1</a:Value></a:KeyValueOfstringstring>
				<a:KeyValueOfstringstring><a:Key>param2</a:Key><a:Value>value2</a:Value></a:KeyValueOfstringstring>
			</GetParametersResult>
		</GetParametersResponse>`))
	result, err := c.GetParameters(ctx(), true)
	if err != nil {
		t.Fatalf("GetParameters: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 params, got %d", len(result))
	}
	if result["param1"] != "value1" {
		t.Errorf("param1: expected value1, got %s", result["param1"])
	}
	if result["param2"] != "value2" {
		t.Errorf("param2: expected value2, got %s", result["param2"])
	}
}

func TestGetParameters_Empty(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<GetParametersResponse xmlns="http://tempuri.org/">
			<GetParametersResult xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/>
		</GetParametersResponse>`))
	result, err := c.GetParameters(ctx(), false)
	if err != nil {
		t.Fatalf("GetParameters: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 params, got %d", len(result))
	}
}

func TestGetStoreSettings(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetStoreSettings",
		soapResp(`<GetStoreSettingsResponse xmlns="http://tempuri.org/">
			<GetStoreSettingsResult xmlns:a="http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol">
				<a:NoAddBonusForAdvertising>true</a:NoAddBonusForAdvertising>
				<a:NoDiscountsForAdvertising>false</a:NoDiscountsForAdvertising>
				<a:NoPayBonusForAdvertising>true</a:NoPayBonusForAdvertising>
				<a:OfflineCheckTime>300</a:OfflineCheckTime>
				<a:OfflineChequeSendCount>10</a:OfflineChequeSendCount>
				<a:OfflineDiscount>true</a:OfflineDiscount>
				<a:OnlineCheckTime>60</a:OnlineCheckTime>
				<a:StoreSettingsID>42</a:StoreSettingsID>
				<a:SyncroTimeout>120</a:SyncroTimeout>
				<a:Timeout>30</a:Timeout>
				<a:UseMapping>false</a:UseMapping>
			</GetStoreSettingsResult>
		</GetStoreSettingsResponse>`))
	result, err := c.GetStoreSettings(ctx(), &StoreConfig{StoreSettingsID: 42})
	if err != nil {
		t.Fatalf("GetStoreSettings: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.StoreSettingsID != 42 {
		t.Errorf("StoreSettingsID: expected 42, got %d", result.StoreSettingsID)
	}
	if !result.NoAddBonusForAdvertising {
		t.Error("NoAddBonusForAdvertising should be true")
	}
	if result.OfflineCheckTime != 300 {
		t.Errorf("OfflineCheckTime: expected 300, got %d", result.OfflineCheckTime)
	}
	if result.OfflineChequeSendCount != 10 {
		t.Errorf("OfflineChequeSendCount: expected 10, got %d", result.OfflineChequeSendCount)
	}
	if !result.OfflineDiscount {
		t.Error("OfflineDiscount should be true")
	}
	if result.Timeout != 30 {
		t.Errorf("Timeout: expected 30, got %d", result.Timeout)
	}
}

func TestGetStoreSettings_NilInput(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<GetStoreSettingsResponse xmlns="http://tempuri.org/">
			<GetStoreSettingsResult xmlns:a="http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol">
				<a:StoreSettingsID>1</a:StoreSettingsID>
				<a:Timeout>10</a:Timeout>
			</GetStoreSettingsResult>
		</GetStoreSettingsResponse>`))
	result, err := c.GetStoreSettings(ctx(), nil)
	if err != nil {
		t.Fatalf("GetStoreSettings: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestRegisterDiscountCard(t *testing.T) {
	c := newTestClientWithValidator(t, "IRSLoyaltyService/RegisterDiscountCard",
		func(body string) {
			if !strings.Contains(body, "12345") {
				t.Error("request should contain card number")
			}
		},
		soapResp(`<RegisterDiscountCardResponse xmlns="http://tempuri.org/"/>`))
	err := c.RegisterDiscountCard(ctx(), "12345")
	if err != nil {
		t.Fatalf("RegisterDiscountCard: %v", err)
	}
}

func TestIsCardValid(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/IsCardValid",
		soapResp(`<IsCardValidResponse xmlns="http://tempuri.org/"><IsCardValidResult>true</IsCardValidResult></IsCardValidResponse>`))
	result, err := c.IsCardValid(ctx(), "CARD001")
	if err != nil {
		t.Fatalf("IsCardValid: %v", err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestIsCardValid_False(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<IsCardValidResponse xmlns="http://tempuri.org/"><IsCardValidResult>false</IsCardValidResult></IsCardValidResponse>`))
	result, err := c.IsCardValid(ctx(), "INVALID")
	if err != nil {
		t.Fatalf("IsCardValid: %v", err)
	}
	if result {
		t.Error("expected false")
	}
}

func TestIsCouponValid(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/IsCouponValid",
		soapResp(`<IsCouponValidResponse xmlns="http://tempuri.org/"><IsCouponValidResult>true</IsCouponValidResult></IsCouponValidResponse>`))
	result, err := c.IsCouponValid(ctx(), "COUPON001")
	if err != nil {
		t.Fatalf("IsCouponValid: %v", err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestGetVerifyCode(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetVerifyCode",
		soapResp(`<GetVerifyCodeResponse xmlns="http://tempuri.org/"><GetVerifyCodeResult>ABC123</GetVerifyCodeResult></GetVerifyCodeResponse>`))
	result, err := c.GetVerifyCode(ctx(), "CARD001")
	if err != nil {
		t.Fatalf("GetVerifyCode: %v", err)
	}
	if result != "ABC123" {
		t.Errorf("expected ABC123, got %s", result)
	}
}

func TestGetCardBalance(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetCardBalance",
		soapResp(`<GetCardBalanceResponse xmlns="http://tempuri.org/"><GetCardBalanceResult>1500.50</GetCardBalanceResult></GetCardBalanceResponse>`))
	result, err := c.GetCardBalance(ctx(), "CARD001")
	if err != nil {
		t.Fatalf("GetCardBalance: %v", err)
	}
	if result != "1500.50" {
		t.Errorf("expected 1500.50, got %s", result)
	}
}

func TestGetCardDiscountAmount(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetCardDiscountAmount",
		soapResp(`<GetCardDiscountAmountResponse xmlns="http://tempuri.org/"><GetCardDiscountAmountResult>125.75</GetCardDiscountAmountResult></GetCardDiscountAmountResponse>`))
	result, err := c.GetCardDiscountAmount(ctx(), "CARD001", "<cheque/>")
	if err != nil {
		t.Fatalf("GetCardDiscountAmount: %v", err)
	}
	if result != 125.75 {
		t.Errorf("expected 125.75, got %f", result)
	}
}

func TestGetCardDiscountAmount_Zero(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<GetCardDiscountAmountResponse xmlns="http://tempuri.org/"><GetCardDiscountAmountResult>0</GetCardDiscountAmountResult></GetCardDiscountAmountResponse>`))
	result, err := c.GetCardDiscountAmount(ctx(), "CARD001", "<cheque/>")
	if err != nil {
		t.Fatalf("GetCardDiscountAmount: %v", err)
	}
	if result != 0 {
		t.Errorf("expected 0, got %f", result)
	}
}

func TestGetCardDiscountAmountString(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetCardDiscountAmountString",
		soapResp(`<GetCardDiscountAmountStringResponse xmlns="http://tempuri.org/"><GetCardDiscountAmountStringResult>125.75 RUB</GetCardDiscountAmountStringResult></GetCardDiscountAmountStringResponse>`))
	result, err := c.GetCardDiscountAmountString(ctx(), "CARD001", "<cheque/>")
	if err != nil {
		t.Fatalf("GetCardDiscountAmountString: %v", err)
	}
	if result != "125.75 RUB" {
		t.Errorf("expected '125.75 RUB', got %s", result)
	}
}

func TestGetMessages(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetMessages",
		soapResp(`<GetMessagesResponse xmlns="http://tempuri.org/"><GetMessagesResult>Hello World</GetMessagesResult></GetMessagesResponse>`))
	result, err := c.GetMessages(ctx(), "<cheque/>")
	if err != nil {
		t.Fatalf("GetMessages: %v", err)
	}
	if result != "Hello World" {
		t.Errorf("expected 'Hello World', got %s", result)
	}
}

func TestGetDiscounts(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetDiscounts",
		soapResp(`<GetDiscountsResponse xmlns="http://tempuri.org/"><GetDiscountsResult>discount_data</GetDiscountsResult></GetDiscountsResponse>`))
	result, err := c.GetDiscounts(ctx(), "<cheque/>")
	if err != nil {
		t.Fatalf("GetDiscounts: %v", err)
	}
	if result != "discount_data" {
		t.Errorf("expected discount_data, got %s", result)
	}
}

func TestGetEmail(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetEmail",
		soapResp(`<GetEmailResponse xmlns="http://tempuri.org/"><GetEmailResult>test@example.com</GetEmailResult></GetEmailResponse>`))
	result, err := c.GetEmail(ctx(), "CARD001")
	if err != nil {
		t.Fatalf("GetEmail: %v", err)
	}
	if result != "test@example.com" {
		t.Errorf("expected test@example.com, got %s", result)
	}
}

func TestGetSelfBuyDiscounts(t *testing.T) {
	c := newTestClientWithValidator(t, "IRSLoyaltyService/GetSelfBuyDiscounts",
		func(body string) {
			if !strings.Contains(body, "storeId") {
				t.Error("request should contain storeId")
			}
		},
		soapResp(`<GetSelfBuyDiscountsResponse xmlns="http://tempuri.org/"><GetSelfBuyDiscountsResult>selfbuy_data</GetSelfBuyDiscountsResult></GetSelfBuyDiscountsResponse>`))
	result, err := c.GetSelfBuyDiscounts(ctx(), "<cheque/>", 100)
	if err != nil {
		t.Fatalf("GetSelfBuyDiscounts: %v", err)
	}
	if result != "selfbuy_data" {
		t.Errorf("expected selfbuy_data, got %s", result)
	}
}

func TestAccrual(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/Accrual",
		soapResp(`<AccrualResponse xmlns="http://tempuri.org/"><AccrualResult>accrual_ok</AccrualResult></AccrualResponse>`))
	result, err := c.Accrual(ctx(), "<cheque/>")
	if err != nil {
		t.Fatalf("Accrual: %v", err)
	}
	if result != "accrual_ok" {
		t.Errorf("expected accrual_ok, got %s", result)
	}
}

func TestOfflineAccrual(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/OfflineAccrual",
		soapResp(`<OfflineAccrualResponse xmlns="http://tempuri.org/"><OfflineAccrualResult>true</OfflineAccrualResult></OfflineAccrualResponse>`))
	result, err := c.OfflineAccrual(ctx(), "<cheque/>")
	if err != nil {
		t.Fatalf("OfflineAccrual: %v", err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestRefund(t *testing.T) {
	c := newTestClientWithValidator(t, "IRSLoyaltyService/Refund",
		func(body string) {
			if !strings.Contains(body, "refundRefundCheque") {
				t.Error("request should contain refundRefundCheque")
			}
			if !strings.Contains(body, "chequeID") {
				t.Error("request should contain chequeID")
			}
		},
		soapResp(`<RefundResponse xmlns="http://tempuri.org/"/>`))
	err := c.Refund(ctx(), "<refund/>", 999)
	if err != nil {
		t.Fatalf("Refund: %v", err)
	}
}

func TestSubtractBonus(t *testing.T) {
	c := newTestClientWithValidator(t, "IRSLoyaltyService/SubtractBonus",
		func(body string) {
			if !strings.Contains(body, "discountCardNumber") {
				t.Error("request should contain discountCardNumber")
			}
			if !strings.Contains(body, "amount") {
				t.Error("request should contain amount")
			}
		},
		soapResp(`<SubtractBonusResponse xmlns="http://tempuri.org/"/>`))
	err := c.SubtractBonus(ctx(), "CARD001", 50.25, "<cheque/>")
	if err != nil {
		t.Fatalf("SubtractBonus: %v", err)
	}
}

func TestSubtractBonus45(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/SubtractBonus45",
		soapResp(`<SubtractBonus45Response xmlns="http://tempuri.org/"><SubtractBonus45Result>ok</SubtractBonus45Result></SubtractBonus45Response>`))
	result, err := c.SubtractBonus45(ctx(), "CARD001", 100.00, "<cheque/>")
	if err != nil {
		t.Fatalf("SubtractBonus45: %v", err)
	}
	if result != "ok" {
		t.Errorf("expected ok, got %s", result)
	}
}

func TestCancelSubtractBonus(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/CancelSubtractBonus",
		soapResp(`<CancelSubtractBonusResponse xmlns="http://tempuri.org/"/>`))
	err := c.CancelSubtractBonus(ctx(), "CARD001", 50.25, "<cheque/>")
	if err != nil {
		t.Fatalf("CancelSubtractBonus: %v", err)
	}
}

func TestValidateUser(t *testing.T) {
	c := newTestClientWithValidator(t, "IRSLoyaltyService/ValidateUser",
		func(body string) {
			if !strings.Contains(body, "userName") {
				t.Error("request should contain userName")
			}
			if !strings.Contains(body, "password") {
				t.Error("request should contain password")
			}
		},
		soapResp(`<ValidateUserResponse xmlns="http://tempuri.org/"><ValidateUserResult>true</ValidateUserResult></ValidateUserResponse>`))
	result, err := c.ValidateUser(ctx(), "admin", "secret")
	if err != nil {
		t.Fatalf("ValidateUser: %v", err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestValidateUser_Invalid(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<ValidateUserResponse xmlns="http://tempuri.org/"><ValidateUserResult>false</ValidateUserResult></ValidateUserResponse>`))
	result, err := c.ValidateUser(ctx(), "admin", "wrong")
	if err != nil {
		t.Fatalf("ValidateUser: %v", err)
	}
	if result {
		t.Error("expected false")
	}
}

func TestCheckDiscountCard(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/CheckDiscountCard",
		soapResp(`<CheckDiscountCardResponse xmlns="http://tempuri.org/"><CheckDiscountCardResult>true</CheckDiscountCardResult></CheckDiscountCardResponse>`))
	result, err := c.CheckDiscountCard(ctx(), "CARD001")
	if err != nil {
		t.Fatalf("CheckDiscountCard: %v", err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestValidateUserRole(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/ValidateUserRole",
		soapResp(`<ValidateUserRoleResponse xmlns="http://tempuri.org/"><ValidateUserRoleResult>true</ValidateUserRoleResult></ValidateUserRoleResponse>`))
	result, err := c.ValidateUserRole(ctx(), "admin", "Manager")
	if err != nil {
		t.Fatalf("ValidateUserRole: %v", err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestGetUserRole(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetUserRole",
		soapResp(`<GetUserRoleResponse xmlns="http://tempuri.org/"><GetUserRoleResult>Administrator</GetUserRoleResult></GetUserRoleResponse>`))
	result, err := c.GetUserRole(ctx(), "admin")
	if err != nil {
		t.Fatalf("GetUserRole: %v", err)
	}
	if result != "Administrator" {
		t.Errorf("expected Administrator, got %s", result)
	}
}

func TestActivationPaymentCard(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/ActivationPaymentCard",
		soapResp(`<ActivationPaymentCardResponse xmlns="http://tempuri.org/"><ActivationPaymentCardResult>500.00</ActivationPaymentCardResult></ActivationPaymentCardResponse>`))
	result, err := c.ActivationPaymentCard(ctx(), "PAY001")
	if err != nil {
		t.Fatalf("ActivationPaymentCard: %v", err)
	}
	if result != 500.00 {
		t.Errorf("expected 500.00, got %f", result)
	}
}

func TestCancelActivationPaymentCard(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/CancelActivationPaymentCard",
		soapResp(`<CancelActivationPaymentCardResponse xmlns="http://tempuri.org/"><CancelActivationPaymentCardResult>true</CancelActivationPaymentCardResult></CancelActivationPaymentCardResponse>`))
	result, err := c.CancelActivationPaymentCard(ctx(), "PAY001")
	if err != nil {
		t.Fatalf("CancelActivationPaymentCard: %v", err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestQuerySyncStream(t *testing.T) {
	c := newTestClientWithValidator(t, "IRSLoyaltyService/QuerySyncStream",
		func(body string) {
			if !strings.Contains(body, "TupleOfstringlong") {
				t.Error("request should contain TupleOfstringlong")
			}
		},
		soapResp(`<QuerySyncStreamResponse xmlns="http://tempuri.org/"><QuerySyncStreamResult>task-123</QuerySyncStreamResult></QuerySyncStreamResponse>`))
	data := []TupleOfStringLong{
		{Item1: "table1", Item2: 100},
		{Item1: "table2", Item2: 200},
	}
	result, err := c.QuerySyncStream(ctx(), data)
	if err != nil {
		t.Fatalf("QuerySyncStream: %v", err)
	}
	if result != "task-123" {
		t.Errorf("expected task-123, got %s", result)
	}
}

func TestQuerySyncStream_Empty(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<QuerySyncStreamResponse xmlns="http://tempuri.org/"><QuerySyncStreamResult>task-0</QuerySyncStreamResult></QuerySyncStreamResponse>`))
	result, err := c.QuerySyncStream(ctx(), nil)
	if err != nil {
		t.Fatalf("QuerySyncStream: %v", err)
	}
	if result != "task-0" {
		t.Errorf("expected task-0, got %s", result)
	}
}

func TestGetSyncStream(t *testing.T) {
	streamData := []byte("binary stream data")
	encoded := base64.StdEncoding.EncodeToString(streamData)
	c := newTestClient(t, "IRSLoyaltyService/GetSyncStream",
		soapResp(`<GetSyncStreamResponse xmlns="http://tempuri.org/"><GetSyncStreamResult>`+encoded+`</GetSyncStreamResult></GetSyncStreamResponse>`))
	result, err := c.GetSyncStream(ctx(), "task-123")
	if err != nil {
		t.Fatalf("GetSyncStream: %v", err)
	}
	if string(result) != string(streamData) {
		t.Errorf("expected %q, got %q", streamData, result)
	}
}

func TestIsTaskCompleted(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/IsTaskCompleted",
		soapResp(`<IsTaskCompletedResponse xmlns="http://tempuri.org/"><IsTaskCompletedResult>true</IsTaskCompletedResult></IsTaskCompletedResponse>`))
	result, err := c.IsTaskCompleted(ctx(), "task-123")
	if err != nil {
		t.Fatalf("IsTaskCompleted: %v", err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestIsTaskCompleted_False(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<IsTaskCompletedResponse xmlns="http://tempuri.org/"><IsTaskCompletedResult>false</IsTaskCompletedResult></IsTaskCompletedResponse>`))
	result, err := c.IsTaskCompleted(ctx(), "task-456")
	if err != nil {
		t.Fatalf("IsTaskCompleted: %v", err)
	}
	if result {
		t.Error("expected false")
	}
}

func TestGetUpdateStream(t *testing.T) {
	streamData := []byte("update binary data")
	encoded := base64.StdEncoding.EncodeToString(streamData)
	c := newTestClient(t, "IRSLoyaltyService/GetUpdateStream",
		soapResp(`<GetUpdateStreamResponse xmlns="http://tempuri.org/"><GetUpdateStreamResult>`+encoded+`</GetUpdateStreamResult></GetUpdateStreamResponse>`))
	result, err := c.GetUpdateStream(ctx(), "update.zip")
	if err != nil {
		t.Fatalf("GetUpdateStream: %v", err)
	}
	if string(result) != string(streamData) {
		t.Errorf("expected %q, got %q", streamData, result)
	}
}

func TestUploadReferences(t *testing.T) {
	c := newTestClientWithValidator(t, "IRSLoyaltyService/UploadReferences",
		func(body string) {
			if !strings.Contains(body, "packet") {
				t.Error("request should contain packet")
			}
			if !strings.Contains(body, "stamp") {
				t.Error("request should contain stamp")
			}
		},
		soapResp(`<UploadReferencesResponse xmlns="http://tempuri.org/"/>`))
	err := c.UploadReferences(ctx(), "ref_data", 42)
	if err != nil {
		t.Fatalf("UploadReferences: %v", err)
	}
}

func TestGetReferencesStamp(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetReferencesStamp",
		soapResp(`<GetReferencesStampResponse xmlns="http://tempuri.org/"><GetReferencesStampResult>12345678</GetReferencesStampResult></GetReferencesStampResponse>`))
	result, err := c.GetReferencesStamp(ctx())
	if err != nil {
		t.Fatalf("GetReferencesStamp: %v", err)
	}
	if result != 12345678 {
		t.Errorf("expected 12345678, got %d", result)
	}
}

func TestGetDataPacket(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetDataPacket",
		soapResp(`<GetDataPacketResponse xmlns="http://tempuri.org/"><GetDataPacketResult>packet_data</GetDataPacketResult></GetDataPacketResponse>`))
	result, err := c.GetDataPacket(ctx(), "query_params")
	if err != nil {
		t.Fatalf("GetDataPacket: %v", err)
	}
	if result != "packet_data" {
		t.Errorf("expected packet_data, got %s", result)
	}
}

func TestSendInfoPacket(t *testing.T) {
	c := newTestClientWithValidator(t, "IRSLoyaltyService/SendInfoPacket",
		func(body string) {
			if !strings.Contains(body, "OfflineChequeCount") {
				t.Error("request should contain OfflineChequeCount")
			}
			if !strings.Contains(body, "Version") {
				t.Error("request should contain Version")
			}
		},
		soapResp(`<SendInfoPacketResponse xmlns="http://tempuri.org/"/>`))
	err := c.SendInfoPacket(ctx(), &RSInfoPacket{
		OfflineChequeCount: 5,
		Version:            "3.0",
	})
	if err != nil {
		t.Fatalf("SendInfoPacket: %v", err)
	}
}

func TestSendInfoPacket_Nil(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<SendInfoPacketResponse xmlns="http://tempuri.org/"/>`))
	err := c.SendInfoPacket(ctx(), nil)
	if err != nil {
		t.Fatalf("SendInfoPacket nil: %v", err)
	}
}

func TestGetStatistic(t *testing.T) {
	c := newTestClient(t, "IRSLoyaltyService/GetStatistic",
		soapResp(`<GetStatisticResponse xmlns="http://tempuri.org/">
			<GetStatisticResult xmlns:a="http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol">
				<a:ItemStatistics>
					<a:DailyQuantity>1.5</a:DailyQuantity>
					<a:ItemId>42</a:ItemId>
					<a:MonthlyQuantity>30.0</a:MonthlyQuantity>
					<a:WeeklyQuantity>7.5</a:WeeklyQuantity>
				</a:ItemStatistics>
				<a:ItemStatistics>
					<a:DailyQuantity>2.0</a:DailyQuantity>
					<a:ItemId>43</a:ItemId>
					<a:MonthlyQuantity>60.0</a:MonthlyQuantity>
					<a:WeeklyQuantity>14.0</a:WeeklyQuantity>
				</a:ItemStatistics>
			</GetStatisticResult>
		</GetStatisticResponse>`))
	result, err := c.GetStatistic(ctx(), GetStatisticRequest{
		AccountID:     1,
		ItemIDs:       []int64{42, 43},
		Time:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		StatisticFlag: StatisticDaily,
	})
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
	if result[0].ItemID != 42 {
		t.Errorf("item[0].ItemID: expected 42, got %d", result[0].ItemID)
	}
	if result[0].DailyQuantity != 1.5 {
		t.Errorf("item[0].DailyQuantity: expected 1.5, got %f", result[0].DailyQuantity)
	}
	if result[1].ItemID != 43 {
		t.Errorf("item[1].ItemID: expected 43, got %d", result[1].ItemID)
	}
	if result[1].MonthlyQuantity != 60.0 {
		t.Errorf("item[1].MonthlyQuantity: expected 60.0, got %f", result[1].MonthlyQuantity)
	}
}

func TestGetStatistic_Empty(t *testing.T) {
	c := newTestClient(t, "",
		soapResp(`<GetStatisticResponse xmlns="http://tempuri.org/">
			<GetStatisticResult xmlns:a="http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol"/>
		</GetStatisticResponse>`))
	result, err := c.GetStatistic(ctx(), GetStatisticRequest{
		AccountID:     1,
		StatisticFlag: StatisticNone,
	})
	if err != nil {
		t.Fatalf("GetStatistic: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 items, got %d", len(result))
	}
}

// ---------- Error handling tests ----------

func TestSOAPFault(t *testing.T) {
	c := newTestClient(t, "",
		soapFaultResp("s:Receiver", "Internal server error", "stack trace here"))
	err := c.Ping(ctx())
	if err == nil {
		t.Fatal("expected error")
	}
	fe, ok := IsFaultError(err)
	if !ok {
		t.Fatalf("expected FaultError, got %T: %v", err, err)
	}
	if fe.Code != "s:Receiver" {
		t.Errorf("expected code s:Receiver, got %s", fe.Code)
	}
	if fe.Reason != "Internal server error" {
		t.Errorf("expected reason 'Internal server error', got %s", fe.Reason)
	}
	if fe.Detail != "stack trace here" {
		t.Errorf("expected detail 'stack trace here', got %s", fe.Detail)
	}
}

func TestSOAPFault_NoDetail(t *testing.T) {
	c := newTestClient(t, "",
		soapFaultResp("s:Sender", "Bad request", ""))
	_, err := c.IsOnline(ctx(), "1.0")
	if err == nil {
		t.Fatal("expected error")
	}
	fe, ok := IsFaultError(err)
	if !ok {
		t.Fatalf("expected FaultError, got %T", err)
	}
	if !strings.Contains(fe.Error(), "Bad request") {
		t.Errorf("error message should contain reason: %s", fe.Error())
	}
}

func TestHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("bad gateway"))
	}))
	t.Cleanup(server.Close)
	c := NewClient(server.URL)
	err := c.Ping(ctx())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "502") {
		t.Errorf("expected 502 in error, got: %v", err)
	}
}

func TestInvalidXMLResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")
		_, _ = w.Write([]byte("not xml at all"))
	}))
	t.Cleanup(server.Close)
	c := NewClient(server.URL)
	err := c.Ping(ctx())
	if err == nil {
		t.Fatal("expected error on invalid XML")
	}
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	t.Cleanup(server.Close)
	c := NewClient(server.URL)
	cancelCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := c.Ping(cancelCtx)
	if err == nil {
		t.Fatal("expected error on context cancellation")
	}
}

func TestServerNotReachable(t *testing.T) {
	c := NewClient("http://127.0.0.1:1") // port 1 - unlikely to be open
	err := c.Ping(ctx())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}

// Note: SOAP envelope and XML type tests are in internal/soap/soap_test.go

// ---------- Client options tests ----------

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("https://example.com/service")
	if c.url != "https://example.com/service" {
		t.Errorf("expected url, got %s", c.url)
	}
	if c.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	c := NewClient("https://example.com/service", WithHTTPClient(custom))
	if c.httpClient != custom {
		t.Error("expected custom http client")
	}
}

func TestNewClient_WithTimeout(t *testing.T) {
	c := NewClient("https://example.com/service", WithTimeout(10*time.Second))
	if c.httpClient.Timeout != 10*time.Second {
		t.Errorf("expected 10s timeout, got %v", c.httpClient.Timeout)
	}
}

// ---------- FaultError tests ----------

func TestFaultError_Error(t *testing.T) {
	fe := &FaultError{Code: "s:Receiver", Reason: "Error occurred"}
	s := fe.Error()
	if !strings.Contains(s, "s:Receiver") || !strings.Contains(s, "Error occurred") {
		t.Errorf("unexpected error string: %s", s)
	}
}

func TestFaultError_ErrorWithDetail(t *testing.T) {
	fe := &FaultError{Code: "s:Receiver", Reason: "Error occurred", Detail: "details"}
	s := fe.Error()
	if !strings.Contains(s, "details") {
		t.Errorf("error string should contain detail: %s", s)
	}
}

func TestIsFaultError_True(t *testing.T) {
	fe := &FaultError{Code: "test"}
	result, ok := IsFaultError(fe)
	if !ok {
		t.Error("expected true")
	}
	if result.Code != "test" {
		t.Error("expected code 'test'")
	}
}

func TestIsFaultError_False(t *testing.T) {
	_, ok := IsFaultError(fmt.Errorf("regular error"))
	if ok {
		t.Error("expected false")
	}
}

// ---------- SOAP 1.2 content-type validation ----------

func TestSOAP12ContentType(t *testing.T) {
	var receivedCT string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedCT = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		_, _ = w.Write([]byte(soapResp(`<PingResponse xmlns="http://tempuri.org/"/>`)))
	}))
	t.Cleanup(server.Close)
	c := NewClient(server.URL)
	_ = c.Ping(ctx())

	if !strings.HasPrefix(receivedCT, "application/soap+xml") {
		t.Errorf("content-type should start with application/soap+xml, got: %s", receivedCT)
	}
	if !strings.Contains(receivedCT, `action="http://tempuri.org/IRSLoyaltyService/Ping"`) {
		t.Errorf("content-type should contain SOAPAction, got: %s", receivedCT)
	}
}

// ---------- WS-Addressing header validation ----------

func TestWSAddressingHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		s := string(body)
		if !strings.Contains(s, "a:Action") {
			t.Error("envelope should contain a:Action header")
		}
		if !strings.Contains(s, "a:To") {
			t.Error("envelope should contain a:To header")
		}
		if !strings.Contains(s, "mustUnderstand") {
			t.Error("envelope should contain mustUnderstand attribute")
		}
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		_, _ = w.Write([]byte(soapResp(`<PingResponse xmlns="http://tempuri.org/"/>`)))
	}))
	t.Cleanup(server.Close)
	c := NewClient(server.URL)
	_ = c.Ping(ctx())
}

// ---------- Integration: SOAP fault on non-void operation ----------

func TestSOAPFault_OnBoolOperation(t *testing.T) {
	c := newTestClient(t, "",
		soapFaultResp("s:Sender", "Card not found", ""))
	_, err := c.IsCardValid(ctx(), "INVALID")
	if err == nil {
		t.Fatal("expected error")
	}
	fe, ok := IsFaultError(err)
	if !ok {
		t.Fatalf("expected FaultError, got %T", err)
	}
	if fe.Reason != "Card not found" {
		t.Errorf("expected 'Card not found', got %s", fe.Reason)
	}
}

func TestSOAPFault_OnStringOperation(t *testing.T) {
	c := newTestClient(t, "",
		soapFaultResp("s:Receiver", "Service unavailable", ""))
	_, err := c.GetCardBalance(ctx(), "CARD001")
	if err == nil {
		t.Fatal("expected error")
	}
	_, ok := IsFaultError(err)
	if !ok {
		t.Fatalf("expected FaultError, got %T", err)
	}
}

func TestSOAPFault_OnDecimalOperation(t *testing.T) {
	c := newTestClient(t, "",
		soapFaultResp("s:Receiver", "Calculation error", ""))
	_, err := c.GetCardDiscountAmount(ctx(), "CARD001", "<cheque/>")
	if err == nil {
		t.Fatal("expected error")
	}
	_, ok := IsFaultError(err)
	if !ok {
		t.Fatalf("expected FaultError, got %T", err)
	}
}

// ---------- SOAP 500 with fault body ----------

func TestHTTP500WithSOAPFault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(soapFaultResp("s:Receiver", "Internal error", "trace")))
	}))
	t.Cleanup(server.Close)
	c := NewClient(server.URL)
	err := c.Ping(ctx())
	if err == nil {
		t.Fatal("expected error")
	}
	fe, ok := IsFaultError(err)
	if !ok {
		t.Fatalf("expected FaultError on 500, got %T: %v", err, err)
	}
	if fe.Reason != "Internal error" {
		t.Errorf("expected 'Internal error', got %s", fe.Reason)
	}
}

