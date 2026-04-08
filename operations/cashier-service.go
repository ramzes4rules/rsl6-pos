package operations

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"

	"github.com/ramzes4rules/rsl6-pos/cheque"
	"github.com/ramzes4rules/rsl6-pos/client"
)

// CashierService encapsulates high-level POS cashier operations.
// It uses [client.Service] to communicate with the loyalty service
// and the [cheque] package to manipulate receipt data.
type CashierService struct {
	client client.Service
	log    *slog.Logger
}

// NewCashierService creates a new CashierService.
//   - svc: an implementation of [client.Service] (e.g. *client.Client).
//   - logger: optional *slog.Logger; if nil, slog.Default() is used.
func NewCashierService(svc client.Service, logger *slog.Logger) *CashierService {
	if logger == nil {
		logger = slog.Default()
	}
	return &CashierService{client: svc, log: logger}
}

// Ping checks connectivity with the loyalty service.
func (s *CashierService) Ping(ctx context.Context) error {
	return s.client.Ping(ctx)
}

// IsOnline checks whether the store module is online.
func (s *CashierService) IsOnline(ctx context.Context, version string) (bool, error) {
	return s.client.IsOnline(ctx, version)
}

// IsCardValid checks if a discount card number is valid.
func (s *CashierService) IsCardValid(ctx context.Context, number string) (bool, error) {
	return s.client.IsCardValid(ctx, number)
}

// IsCouponValid checks if a coupon number is valid.
func (s *CashierService) IsCouponValid(ctx context.Context, coupon string) (bool, error) {
	return s.client.IsCouponValid(ctx, coupon)
}

// GetCardBalance retrieves card balance information.
// The raw XML response is parsed into [cheque.CardBalanceResult].
func (s *CashierService) GetCardBalance(ctx context.Context, cardNumber string) (*cheque.CardBalanceResult, error) {
	raw, err := s.client.GetCardBalance(ctx, cardNumber)
	if err != nil {
		return nil, fmt.Errorf("get card balance: %w", err)
	}
	var result cheque.CardBalanceResult
	if err := xml.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("parse card balance response: %w", err)
	}
	return &result, nil
}

// GetMessages retrieves loyalty messages for the receipt, parses them,
// and applies them to the receipt.
func (s *CashierService) GetMessages(ctx context.Context, receipt *cheque.Cheque) (cheque.LoyaltyMessages, error) {
	receipt.DeleteMessages()

	raw, err := s.client.GetMessages(ctx, receipt.SerializeToXml(false))
	if err != nil {
		return cheque.LoyaltyMessages{}, fmt.Errorf("get messages: %w", err)
	}

	var messages cheque.LoyaltyMessages
	if err := xml.Unmarshal([]byte(raw), &messages); err != nil {
		return cheque.LoyaltyMessages{}, fmt.Errorf("parse messages response: %w", err)
	}

	s.log.Debug("messages fetched", "count", len(messages.Messages))
	receipt.AddMessagesToReceipt(messages)
	return messages, nil
}

// GetDiscounts retrieves loyalty discounts for the receipt, parses them,
// and applies them to the receipt.
func (s *CashierService) GetDiscounts(ctx context.Context, receipt *cheque.Cheque) (cheque.LoyaltyDiscounts, error) {
	receipt.DeleteLoyaltyDiscounts()

	raw, err := s.client.GetDiscounts(ctx, receipt.SerializeToXml(false))
	if err != nil {
		return cheque.LoyaltyDiscounts{}, fmt.Errorf("get discounts: %w", err)
	}

	var discounts cheque.LoyaltyDiscounts
	if err := xml.Unmarshal([]byte(raw), &discounts); err != nil {
		return cheque.LoyaltyDiscounts{}, fmt.Errorf("parse discounts response: %w", err)
	}

	s.log.Debug("discounts fetched", "lines", len(discounts.ChequeLines.ChequeLine))
	receipt.AddLoyaltyDiscounts(discounts)
	return discounts, nil
}

// GetCardDiscountAmount retrieves the available bonus payment amount for a card.
func (s *CashierService) GetCardDiscountAmount(ctx context.Context, cardNumber string, receipt *cheque.Cheque) (float64, error) {
	amount, err := s.client.GetCardDiscountAmount(ctx, cardNumber, receipt.SerializeToXml(false))
	if err != nil {
		return 0, fmt.Errorf("get card discount amount: %w", err)
	}
	return amount, nil
}

// SubtractBonus45 performs bonus subtraction with line-level detail
// and applies the result to the receipt.
func (s *CashierService) SubtractBonus45(ctx context.Context, receipt *cheque.Cheque, cardNumber string, amount float64) (cheque.Subtractions, error) {
	raw, err := s.client.SubtractBonus45(ctx, cardNumber, amount, receipt.SerializeToXml(false))
	if err != nil {
		return nil, fmt.Errorf("subtract bonus 4.5: %w", err)
	}

	var result struct {
		ChequeLines cheque.Subtractions `xml:"ChequeLine"`
	}
	if err := xml.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("parse subtraction response: %w", err)
	}

	if len(result.ChequeLines) > 0 {
		receipt.ApplySubtraction(result.ChequeLines, cardNumber)
	}
	return result.ChequeLines, nil
}

// CancelSubtractBonus cancels a previously performed bonus payment.
func (s *CashierService) CancelSubtractBonus(ctx context.Context, receipt *cheque.Cheque, cardNumber string, amount float64) error {
	if err := s.client.CancelSubtractBonus(ctx, cardNumber, amount, receipt.SerializeToXml(false)); err != nil {
		return fmt.Errorf("cancel subtract bonus: %w", err)
	}
	receipt.DeleteBonusPayment()
	return nil
}

// CancelLoyaltyDiscount clears all loyalty discounts from a receipt.
func (s *CashierService) CancelLoyaltyDiscount(receipt *cheque.Cheque) {
	receipt.DeleteLoyaltyDiscounts()
}

// SubTotal performs the subtotal operation: fetches messages and discounts.
func (s *CashierService) SubTotal(ctx context.Context, receipt *cheque.Cheque) error {
	s.log.Debug("fetching receipt messages")
	if _, err := s.GetMessages(ctx, receipt); err != nil {
		return fmt.Errorf("subtotal messages: %w", err)
	}

	s.log.Debug("fetching receipt discounts")
	if _, err := s.GetDiscounts(ctx, receipt); err != nil {
		return fmt.Errorf("subtotal discounts: %w", err)
	}
	return nil
}

// ChequeClose closes a receipt and performs accrual.
// Returns the slip message from the loyalty service (may be empty).
func (s *CashierService) ChequeClose(ctx context.Context, receipt *cheque.Cheque) (string, error) {
	receipt.CloseReceipt()
	slip, err := s.client.Accrual(ctx, receipt.SerializeToXml(false))
	if err != nil {
		return "", fmt.Errorf("accrual: %w", err)
	}
	return slip, nil
}

// ChequeCancel cancels a receipt. If expert is false, any existing bonus
// subtraction is automatically cancelled first.
func (s *CashierService) ChequeCancel(ctx context.Context, receipt *cheque.Cheque, expert bool) (string, error) {
	s.log.Debug("cancelling receipt", "expert", expert)

	if !expert && receipt.GetSubtractedBonus() != 0 {
		cardNo := receipt.GetLoyaltyCardNumber(0)
		if cardNo != nil {
			if err := s.CancelSubtractBonus(ctx, receipt, *cardNo, receipt.GetSubtractedBonus()); err != nil {
				s.log.Error("failed to cancel bonus subtraction", "error", err)
			}
		}
	}

	receipt.CancelReceipt()

	slip, err := s.client.Accrual(ctx, receipt.SerializeToXml(false))
	if err != nil {
		return "", fmt.Errorf("accrual on cancel: %w", err)
	}
	return slip, nil
}

// PayWithGiftCard performs a gift card payment.
func (s *CashierService) PayWithGiftCard(ctx context.Context, receipt *cheque.Cheque, cardNumber string, amount float64) error {
	if err := s.client.SubtractBonus(ctx, cardNumber, amount, receipt.SerializeToXml(false)); err != nil {
		return fmt.Errorf("gift card payment: %w", err)
	}
	receipt.ApplyGiftCardPayment(cardNumber, amount)
	return nil
}

// CheckConnection performs a full connectivity check: ping, online status, and auth.
func (s *CashierService) CheckConnection(ctx context.Context, version string) error {
	if err := s.client.Ping(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	online, err := s.client.IsOnline(ctx, version)
	if err != nil {
		return fmt.Errorf("is online check failed: %w", err)
	}
	if !online {
		return fmt.Errorf("store module is offline")
	}

	_, err = s.client.IsCouponValid(ctx, "123")
	if err != nil {
		if fe, ok := client.IsFaultError(err); ok && fe.Reason == "Coupon not found" {
			return nil // expected error — auth works
		}
		return fmt.Errorf("auth check failed: %w", err)
	}
	return nil
}

// ActivatePaymentCard activates a payment card and returns the balance.
func (s *CashierService) ActivatePaymentCard(ctx context.Context, cardNumber string) (float64, error) {
	balance, err := s.client.ActivationPaymentCard(ctx, cardNumber)
	if err != nil {
		return 0, fmt.Errorf("activate payment card: %w", err)
	}
	return balance, nil
}

// SerializeReceipt returns the XML representation of a receipt.
func (s *CashierService) SerializeReceipt(receipt *cheque.Cheque, indent bool) string {
	return receipt.SerializeToXml(indent)
}
