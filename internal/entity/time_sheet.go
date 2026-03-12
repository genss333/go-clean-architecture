package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type TimeSheet struct {
	TimeSheetID int
	Employee    Employee
	HourWorked  decimal.Decimal
	WorkedDate  time.Time
}
