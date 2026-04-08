package rsl6_pos_operation

type PosOperation interface {
	GetDiscount()
	AddLoyaltyDiscount()
	AddPosDiscount()
	ClearDiscount()
	GetMessages()
	AddMessagesToCheck()

}
