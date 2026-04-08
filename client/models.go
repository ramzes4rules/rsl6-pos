package client

import "time"

// StoreConfig represents store configuration settings from the loyalty service.
type StoreConfig struct {
	NoAddBonusForAdvertising  bool
	NoDiscountsForAdvertising bool
	NoPayBonusForAdvertising  bool
	OfflineCheckTime          int64
	OfflineChequeSendCount    int32
	OfflineDiscount           bool
	OnlineCheckTime           int64
	StoreSettingsID           int64
	SyncroTimeout             int64
	Timeout                   int64
	UseMapping                bool
}

// RSInfoPacket represents an information packet sent to the service.
type RSInfoPacket struct {
	OfflineChequeCount int64
	Version            string
}

// StatisticPeriodicityType defines periodicity for statistics queries.
type StatisticPeriodicityType string

const (
	StatisticNone    StatisticPeriodicityType = "None"
	StatisticDaily   StatisticPeriodicityType = "Daily"
	StatisticWeekly  StatisticPeriodicityType = "Weekly"
	StatisticMonthly StatisticPeriodicityType = "Monthly"
)

// ItemStatistics represents statistics for a single item.
type ItemStatistics struct {
	DailyQuantity   float64
	ItemID          int64
	MonthlyQuantity float64
	WeeklyQuantity  float64
}

// TupleOfStringLong represents a tuple of string and int64.
type TupleOfStringLong struct {
	Item1 string
	Item2 int64
}

// GetStatisticRequest holds parameters for the GetStatistic operation.
type GetStatisticRequest struct {
	AccountID     int64
	ItemIDs       []int64
	Time          time.Time
	StatisticFlag StatisticPeriodicityType
}
