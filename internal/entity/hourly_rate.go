package entity

import "github.com/shopspring/decimal"

type HourlyRate struct {
	HourlyRateID int
	Amount       decimal.Decimal
}
