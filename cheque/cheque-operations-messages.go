package cheque

import "github.com/ramzes4rules/rsl6-pos/models"

// AddMessagesToReceipt append messages to receipt
func (cheque *models.Cheque) AddMessagesToReceipt(messages models.LoyaltyMessages) {
	if len(messages.Messages) > 0 {
		cheque.Messages = new(models.Messages)
		for _, message := range messages.Messages {
			cheque.Messages.Messages = append(cheque.Messages.Messages, models.Message{
				MessageID: message.MessageID,
				Device:    message.Device,
				Body:      message.Body,
			})
		}
	}
}

// DeleteMessages removes all messages from receipt
func (cheque *models.Cheque) DeleteMessages() {
	cheque.Messages = nil
}
