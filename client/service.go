package rslpos

import "context"

// Service defines the RSLoyaltyService client interface.
// All methods accept a context for cancellation and timeout control.
type Service interface {
	// Ping checks service availability.
	Ping(ctx context.Context) error

	// IsOnline checks if the service is online for the given version.
	IsOnline(ctx context.Context, version string) (bool, error)

	// GetParameters returns service parameters as key-value pairs.
	GetParameters(ctx context.Context, isStore bool) (map[string]string, error)

	// GetStoreSettings retrieves store settings, optionally sending current settings.
	GetStoreSettings(ctx context.Context, currentSettings *StoreConfig) (*StoreConfig, error)

	// RegisterDiscountCard registers a new discount card.
	RegisterDiscountCard(ctx context.Context, discountCardNumber string) error

	// IsCardValid checks if a discount card number is valid.
	IsCardValid(ctx context.Context, discountCardNumber string) (bool, error)

	// IsCouponValid checks if a coupon number is valid.
	IsCouponValid(ctx context.Context, couponNumber string) (bool, error)

	// GetVerifyCode retrieves a verification code for a discount card.
	GetVerifyCode(ctx context.Context, discountCardNumber string) (string, error)

	// GetCardBalance retrieves the balance for a discount card.
	GetCardBalance(ctx context.Context, discountCardNumber string) (string, error)

	// GetCardDiscountAmount calculates the discount amount for a cheque.
	GetCardDiscountAmount(ctx context.Context, discountCardNumber, cheque string) (float64, error)

	// GetCardDiscountAmountString calculates the discount amount as a string.
	GetCardDiscountAmountString(ctx context.Context, discountCardNumber, cheque string) (string, error)

	// GetMessages retrieves messages for a cheque.
	GetMessages(ctx context.Context, cheque string) (string, error)

	// GetDiscounts retrieves discounts for a cheque.
	GetDiscounts(ctx context.Context, cheque string) (string, error)

	// GetEmail retrieves the email for a discount number.
	GetEmail(ctx context.Context, discountNumber string) (string, error)

	// GetSelfBuyDiscounts retrieves self-buy discounts for a cheque and store.
	GetSelfBuyDiscounts(ctx context.Context, cheque string, storeID int64) (string, error)

	// Accrual performs bonus accrual for a cheque.
	Accrual(ctx context.Context, cheque string) (string, error)

	// OfflineAccrual performs offline bonus accrual.
	OfflineAccrual(ctx context.Context, cheque string) (bool, error)

	// Refund processes a refund for a cheque.
	Refund(ctx context.Context, refundCheque string, chequeID int64) error

	// SubtractBonus subtracts bonus from a card.
	SubtractBonus(ctx context.Context, discountCardNumber string, amount float64, cheque string) error

	// SubtractBonus45 subtracts bonus (v4.5) and returns result.
	SubtractBonus45(ctx context.Context, discountCardNumber string, amount float64, cheque string) (string, error)

	// CancelSubtractBonus cancels a previous bonus subtraction.
	CancelSubtractBonus(ctx context.Context, discountCardNumber string, amount float64, cheque string) error

	// ValidateUser validates user credentials.
	ValidateUser(ctx context.Context, userName, password string) (bool, error)

	// CheckDiscountCard checks a discount card.
	CheckDiscountCard(ctx context.Context, discountCard string) (bool, error)

	// ValidateUserRole validates if a user has a specific role.
	ValidateUserRole(ctx context.Context, userName, roleName string) (bool, error)

	// GetUserRole retrieves the role for a user.
	GetUserRole(ctx context.Context, userName string) (string, error)

	// ActivationPaymentCard activates a payment card and returns the amount.
	ActivationPaymentCard(ctx context.Context, discountCard string) (float64, error)

	// CancelActivationPaymentCard cancels payment card activation.
	CancelActivationPaymentCard(ctx context.Context, discountCard string) (bool, error)

	// QuerySyncStream queries a sync stream and returns a task ID.
	QuerySyncStream(ctx context.Context, data []TupleOfStringLong) (string, error)

	// GetSyncStream retrieves sync stream data by task ID.
	GetSyncStream(ctx context.Context, taskID string) ([]byte, error)

	// IsTaskCompleted checks if an async task is completed.
	IsTaskCompleted(ctx context.Context, taskID string) (bool, error)

	// GetUpdateStream retrieves an update stream by filename.
	GetUpdateStream(ctx context.Context, filename string) ([]byte, error)

	// UploadReferences uploads reference data.
	UploadReferences(ctx context.Context, packet string, stamp int64) error

	// GetReferencesStamp retrieves the current references stamp.
	GetReferencesStamp(ctx context.Context) (int64, error)

	// GetDataPacket retrieves a data packet.
	GetDataPacket(ctx context.Context, paramPacket string) (string, error)

	// SendInfoPacket sends an information packet.
	SendInfoPacket(ctx context.Context, packet *RSInfoPacket) error

	// GetStatistic retrieves item statistics.
	GetStatistic(ctx context.Context, req GetStatisticRequest) ([]ItemStatistics, error)
}

