package entity_test

import (
	"testing"

	"github.com/genss333/go-clean-architecture/internal/entity"
	testhelper "github.com/genss333/go-clean-architecture/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCalGrossPay(t *testing.T) {
	rows := testhelper.LoadCSV(t, "cal_gross_pay.csv")

	for _, row := range rows {
		t.Run(row["name"], func(t *testing.T) {
			emp := entity.Employee{
				EmployeeID: 1,
				FullName:   row["employee"],
				Department: entity.Department{DepartmentID: 1, Name: row["department"]},
				HourlyRate: entity.HourlyRate{HourlyRateID: 1, Amount: testhelper.DecimalFrom(t, row["hourly_rate"])},
			}

			payStub := entity.PayStubs{
				PayStubID: 1,
				Employee:  emp,
			}

			payStub.CalGrossPay(testhelper.DecimalFrom(t, row["hour_worked"]))

			expected := testhelper.DecimalFrom(t, row["expected_gross"])
			assert.True(t, expected.Equal(payStub.GrossPay),
				"expected GrossPay %s but got %s", row["expected_gross"], payStub.GrossPay.String())
		})
	}
}

func TestCalNetPay(t *testing.T) {
	rows := testhelper.LoadCSV(t, "cal_net_pay.csv")

	for _, row := range rows {
		t.Run(row["name"], func(t *testing.T) {
			emp := entity.Employee{
				EmployeeID: 1,
				FullName:   row["employee"],
				Department: entity.Department{DepartmentID: 1, Name: row["department"]},
				HourlyRate: entity.HourlyRate{HourlyRateID: 1, Amount: testhelper.DecimalFrom(t, row["hourly_rate"])},
			}

			payStub := entity.PayStubs{
				PayStubID: 1,
				Employee:  emp,
				TaxAmount: testhelper.DecimalFrom(t, row["tax_rate"]),
			}

			payStub.CalGrossPay(testhelper.DecimalFrom(t, row["hour_worked"]))
			payStub.CalNetPay()

			expectedGross := testhelper.DecimalFrom(t, row["expected_gross"])
			expectedNet := testhelper.DecimalFrom(t, row["expected_net"])

			assert.True(t, expectedGross.Equal(payStub.GrossPay),
				"expected GrossPay %s but got %s", row["expected_gross"], payStub.GrossPay.String())
			assert.True(t, expectedNet.Equal(payStub.NetPay),
				"expected NetPay %s but got %s", row["expected_net"], payStub.NetPay.String())
		})
	}
}

func newPayStub() entity.PayStubs {
	return entity.PayStubs{
		PayStubID: 1,
		Employee: entity.Employee{
			EmployeeID: 1,
			FullName:   "John Doe",
			Department: entity.Department{DepartmentID: 1, Name: "Engineering"},
			HourlyRate: entity.HourlyRate{HourlyRateID: 1, Amount: decimal.NewFromInt(500)},
		},
		TaxAmount: decimal.NewFromFloat(0.07),
	}
}

// BenchmarkCalGrossPay วัดประสิทธิภาพการคำนวณ gross pay
func BenchmarkCalGrossPay(b *testing.B) {
	hourWorked := decimal.NewFromInt(160)
	b.ResetTimer()
	for i := range b.N {
		_ = i
		p := newPayStub()
		p.CalGrossPay(hourWorked)
	}
}

// BenchmarkCalNetPay วัดประสิทธิภาพการคำนวณ net pay (gross + หักภาษี)
func BenchmarkCalNetPay(b *testing.B) {
	hourWorked := decimal.NewFromInt(160)
	b.ResetTimer()
	for i := range b.N {
		_ = i
		p := newPayStub()
		p.CalGrossPay(hourWorked)
		p.CalNetPay()
	}
}
