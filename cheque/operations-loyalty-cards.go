package cheque

import "fmt"

// AddLoyaltyCards sets the loyalty cards on the receipt from the provided list.
func (cheque *Cheque) AddLoyaltyCards(cards []string) {
	if len(cards) > 0 {
		cheque.DiscountCard = []DiscountCard{}
		for _, card := range cards {
			if card != "" {
				cheque.DiscountCard = append(cheque.DiscountCard, DiscountCard{DiscountCardNo: card, BonusCard: true})
			}
		}
	} else {
		cheque.DiscountCard = nil
	}
}

// GetLoyaltyCards returns the list of loyalty cards applied to the receipt
func (cheque *Cheque) GetLoyaltyCards() (cards string) {
	for _, card := range cheque.DiscountCard {
		cards += fmt.Sprintf("%s,", card.DiscountCardNo)
	}
	if cards != "" {
		cards = cards[:len(cards)-1]
	}
	return cards
}

// GetLoyaltyCardNumber returns specified loyalty card number
func (cheque *Cheque) GetLoyaltyCardNumber(number int) *string {
	if len(cheque.DiscountCard) >= number {
		return &cheque.DiscountCard[number].DiscountCardNo
	}
	return nil
}
