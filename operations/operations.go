package operations

import (
	"fmt"
	//"schweik2/internal/pos-client/cheque"
	//cheque_discount "schweik2/internal/pos-client/cheque-discount"
	//
	//"github.com/rs/zerolog/log"
	"github.com/ramzes4rules/rsl6-pos/cheque"
	"github.com/ramzes4rules/rsl6-pos/client"
	"github.com/ramzes4rules/rsl6-pos/models"
)

// CashierService encapsulates all POS cashier operations.
// It uses POSClient interface to communicate with the loyalty service
// and is completely UI-independent.
type CashierService struct {
	NewClient func() POSClient
}

// POSClient is a minimal interface matching the existing pos_client.RSLoyaltyPOS5 methods.
// This avoids circular imports with the domain package.
type POSClient interface {
	Ping() error
	IsOnline() (*bool, error)
	IsCardValid(number string) (*bool, error)
	IsCouponValid(couponNumber string) (*bool, error)
	GetCardBalance(number string) (*models.CardBalanceResult, error)
	GetCardDiscountAmount(discountCardNumber string, cheque string) (float64, error)
	GetCardDiscountAmountString(discountCardNumber string, cheque string) (string, error)
	GetMessages(currentCheque string) (cheque.IncomeMessages, error)
	GetDiscounts(cheque string) (models.ChequeDiscounts, error)
	Accrual(cheque string) (*string, error)
	SubtractBonus(discountCardNumber string, amount float64, cheque string) (*string, error)
	SubtractBonus45(discountCardNumber string, amount float64, cheque string) (SubtractedResult, error)
	CancelSubtractBonus(discountCardNumber string, amount float64, cheque string) error
	ActivationPaymentCard(number string) (*float32, error)
}

// SubtractedResult represents the result of SubtractBonus45.
type SubtractedResult struct {
	ChequeLines []cheque.SubtractedChequeLine
}

// NewCashierService creates a new CashierService with the given client factory.
func NewCashierService(clientFactory func() POSClient) *CashierService {
	return &CashierService{NewClient: clientFactory}
}

// Ping checks connectivity with the loyalty service.
func (s *CashierService) Ping() string {
	client := s.NewClient()
	if err := client.Ping(); err != nil {
		return err.Error()
	}
	msg := "Ping success"
	return msg
}

// IsOnLine checks whether the store module is online.
func (s *CashierService) IsOnLine() string {
	client := s.NewClient()
	online, err := client.IsOnline()
	if err != nil {
		return err.Error()
	}
	msg := fmt.Sprintf("Store is online: %t", *online)
	return msg
}

// IsCardValid checks if a discount card is valid.
func (s *CashierService) IsCardValid(number string) string {
	client := s.NewClient()
	valid, err := client.IsCardValid(number)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Is card valid: %t", *valid)
}

// IsCouponValid checks if a coupon is valid.
func (s *CashierService) IsCouponValid(coupon string) string {
	client := s.NewClient()
	valid, err := client.IsCouponValid(coupon)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Coupon %s valid: %t", coupon, *valid)
}

// GetCardBalance retrieves card balance info.
func (s *CashierService) GetCardBalance(number string, gift bool) (*string, error) {
	client := s.NewClient()
	balance, err := client.GetCardBalance(number)
	if err != nil {
		return nil, err
	}
	log.Debug().Str("body", balance.Msg.Body).Msg("GetCardBalance request was performed:")
	if !gift {
		msg := balance.Msg.Body
		return &msg, nil
	}
	msg := fmt.Sprintf("%9.02f", balance.Balance.Value)
	return &msg, nil
}

// GetMessages retrieves messages from the loyalty system for a receipt.
func (s *CashierService) GetMessages(receipt *cheque.Cheque) (*string, error) {
	log.Debug().Msg("Running get messages operation")

	receipt.ClearMessages()
	client := s.NewClient()
	result, err := client.GetMessages(receipt.GetXml(false))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get messages:")
		return nil, err
	}

	msg := ""
	if len(result.Messages) > 0 {
		log.Debug().Int("amount", len(result.Messages)).Msg("Messages were fetched:")
		receipt.AddMessagesToReceipt(result)
		for _, message := range result.Messages {
			msg += fmt.Sprintf("MessageID: %s, Device: %s, Body: %s\n", message.MessageID, message.Device, message.Body)
		}
	} else {
		msg = "No messages"
	}
	return &msg, nil
}

// GetDiscounts retrieves discounts for a receipt.
func (s *CashierService) GetDiscounts(receipt *cheque.Cheque) (*string, error) {
	log.Debug().Msg("Calling procedure to clear existing discounts in receipt")
	receipt.ClearDiscounts()

	client := s.NewClient()
	log.Debug().Str("cheque", receipt.GetXml(true)).Msg("Calling procedure to get receipt discounts:")
	response, err := client.GetDiscounts(receipt.GetXml(false))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get receipt discounts:")
		return nil, err
	}
	log.Debug().Int("amount", len(response.ChequeLines.ChequeLine)).Msg("Receipt discounts were received:")

	msg := ""
	if len(response.ChequeLines.ChequeLine) > 0 {
		receipt.ApplyDiscounts(response)
		msg = ComposeReceiptDiscountMessage(response)
	} else {
		msg = "No discounts found for this receipt"
	}
	return &msg, nil
}

// ComposeReceiptDiscountMessage formats discount data into a human-readable string.
func ComposeReceiptDiscountMessage(discounts cheque_discount.ChequeDiscounts) string {
	msg := ""
	for _, line := range discounts.ChequeLines.ChequeLine {
		msg += fmt.Sprintf("Line number: %d, total amount: %5.2f\n", line.ChequeLineNo, line.TotalAmount)
		for _, dis := range line.Discounts.Discount {
			msg += fmt.Sprintf("-> DiscountID: %d, amount: %5.2f\n", dis.DiscountID, dis.Amount)
		}
	}
	return msg
}

// GetCardDiscountAmount retrieves the available bonus payment amount.
func (s *CashierService) GetCardDiscountAmount(number string, receipt *cheque.Cheque) (*string, float64, error) {
	client := s.NewClient()
	result, err := client.GetCardDiscountAmount(number, receipt.GetXml(false))
	if err != nil {
		return nil, 0, err
	}
	msg := fmt.Sprintf("Available for payment: %9.2f", result)
	return &msg, result, nil
}

// SubtractBonus45 performs bonus subtraction with line-level detail.
func (s *CashierService) SubtractBonus45(receipt *cheque.Cheque, discountCardNo string, amount float64) (*string, error) {
	client := s.NewClient()
	lines, err := client.SubtractBonus45(discountCardNo, amount, receipt.GetXml(false))
	if err != nil {
		return nil, err
	}

	if len(lines.ChequeLines) > 0 {
		receipt.ApplySubtraction(lines.ChequeLines, discountCardNo)
		msg := "Payment distribution:\n"
		for _, line := range lines.ChequeLines {
			for _, chequeLine := range receipt.ChequeLines.ChequeLines {
				if chequeLine.ChequeLineNo == line.ChequeLineNo {
					msg += fmt.Sprintf("Line: %d, amount: %7.2f\n", line.ChequeLineNo, line.Amount)
				}
			}
		}
		return &msg, nil
	}
	noSubtraction := "Bonus payment is not available"
	return &noSubtraction, nil
}

// CancelSubtractBonus cancels a bonus payment.
func (s *CashierService) CancelSubtractBonus(discountCardNo string, amount float64, receipt *cheque.Cheque) (*string, error) {
	client := s.NewClient()
	if err := client.CancelSubtractBonus(discountCardNo, amount, receipt.GetXml(false)); err != nil {
		return nil, err
	}
	receipt.ClearBonusPayment()
	msg := "Bonus payment successfully cancelled"
	return &msg, nil
}

// CancelLoyaltyDiscount clears all loyalty discounts from a receipt.
func (s *CashierService) CancelLoyaltyDiscount(receipt *cheque.Cheque) {
	receipt.ClearDiscounts()
}

// SubTotal performs the subtotal operation: fetches messages and discounts.
func (s *CashierService) SubTotal(receipt *cheque.Cheque) (string, error) {
	log.Debug().Msg("==> Fetching receipt messages from the loyalty system.")
	msg, err := s.GetMessages(receipt)
	if err != nil {
		log.Error().Err(err).Msg("==> Failed to fetch receipt messages:")
		return "", err
	}
	log.Debug().Str("message", *msg).Msg("==> Receipt message was fetched:")

	log.Debug().Msg("==> Fetching loyalty discounts from the loyalty system")
	discounts, err := s.GetDiscounts(receipt)
	if err != nil {
		log.Error().Err(err).Msg("==> Failed to fetch loyalty discount:")
		return "", err
	}
	log.Debug().Str("discount", *discounts).Msg("==> Loyalty discount was fetched:")

	response := fmt.Sprintf("Messages:\n%s\n", *msg)
	response += fmt.Sprintf("Discounts:\n%s\n", *discounts)
	return response, nil
}

// ChequeClose closes a receipt (accrual).
func (s *CashierService) ChequeClose(receipt *cheque.Cheque) (*string, error) {
	receipt.CloseReceipt()
	client := s.NewClient()
	slip, err := client.Accrual(receipt.GetXml(false))
	if err != nil {
		return nil, err
	}
	msg := "Success!\n"
	if slip != nil {
		msg += fmt.Sprintf("Slip message:\n%s", *slip)
	}
	return &msg, nil
}

// ChequeCancel cancels a receipt.
func (s *CashierService) ChequeCancel(receipt *cheque.Cheque, expert bool) (*string, error) {
	log.Debug().Bool("expert", expert).Msg("Execute receipt cancelling procedure:")

	if !expert {
		if receipt.GetSubtractedBonus() != 0 {
			if _, err := s.CancelSubtractBonus(receipt.GetLoyaltyCardNumber(0), receipt.GetSubtractedBonus(), receipt); err != nil {
				log.Error().Err(err).Msg("Failed to cancel bonus subtraction:")
			}
		}
	}

	receipt.CancelReceipt()

	client := s.NewClient()
	slip, err := client.Accrual(receipt.GetXml(false))
	if err != nil {
		return nil, err
	}
	msg := "Success!\n"
	if slip != nil {
		msg += fmt.Sprintf("Slip message:\n%s", *slip)
	}
	return &msg, nil
}

// PayWithCard performs a gift card payment.
func (s *CashierService) PayWithCard(discountCardNo string, amount float64, receipt *cheque.Cheque) (*string, error) {
	client := s.NewClient()
	message, err := client.SubtractBonus(discountCardNo, amount, receipt.GetXml(false))
	if err != nil {
		return nil, err
	}
	receipt.ApplyGiftCardPayment(discountCardNo, amount)
	return message, nil
}

// CheckConnection performs a full connectivity check: ping, online status, and auth.
func (s *CashierService) CheckConnection() string {
	client := s.NewClient()

	if err := client.Ping(); err != nil {
		return "Error: no connection with the store module."
	}

	online, err := client.IsOnline()
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	if !*online {
		return "Error: store module is offline"
	}

	_, err = client.IsCouponValid("123")
	if err != nil && err.Error() != "Coupon not found (99503)" {
		return fmt.Sprintf("Error: %s", err.Error())
	}

	return "Connection is OK!\nService is available,\nstatus is online,\nlogin and password are correct."
}

// GetXml returns the XML representation of a receipt.
func (s *CashierService) GetXml(receipt *cheque.Cheque, indent bool) string {
	return receipt.GetXml(indent)
}

// ActivationPaymentCard activates a payment card.
func (s *CashierService) ActivationPaymentCard(receipt *cheque.Cheque, number string) (string, *float32) {
	client := s.NewClient()
	balance, err := client.ActivationPaymentCard(number)
	if err != nil {
		return fmt.Sprintf("Failed to activate: %s", err), nil
	}
	return fmt.Sprintf("Card activated, balance: %f", *balance), balance
}
