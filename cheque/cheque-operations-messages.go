package cheque

// AddMessagesToReceipt appends loyalty messages to the receipt.
func (cheque *Cheque) AddMessagesToReceipt(messages LoyaltyMessages) {
	if len(messages.Messages) > 0 {
		cheque.Messages = new(Messages)
		for _, m := range messages.Messages {
			cheque.Messages.Messages = append(cheque.Messages.Messages, Message{
				MessageID: m.MessageID,
				Device:    m.Device,
				Body:      m.Body,
			})
		}
	}
}

// DeleteMessages removes all messages from the receipt.
func (cheque *Cheque) DeleteMessages() {
	cheque.Messages = nil
}
