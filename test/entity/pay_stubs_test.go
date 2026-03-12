package entity_test

import (
	"testing"

	"github.com/genss333/go-clean-architecture/internal/entity"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCalGrossPay(t *testing.T) {

	hourRate := entity.HourlyRate{
		HourlyRateID: 1,
		Amount:       decimal.NewFromFloat(50.50),
	}
	emp := entity.Employee{
		EmployeeID: 1,
		FullName:   "peerawas.k",
		Department: entity.Department{
			DepartmentID: 1,
			Name:         "Dev",
		},
		HourlyRate: hourRate,
	}

	timeSheet := entity.TimeSheet{
		HourWorked: decimal.NewFromFloat(8.0),
	}

	e := entity.PayStubs{
		PayStubID: 1,
		Employee:  emp,
	}

	e.CalGrossPay(timeSheet.HourWorked)

	expected := decimal.NewFromInt(404)
	assert.True(t, expected.Equal(e.GrossPay))

}

func TestCalNetPay(t *testing.T) {
	hourRate := entity.HourlyRate{
		HourlyRateID: 1,
		Amount:       decimal.NewFromFloat(50.50),
	}
	emp := entity.Employee{
		EmployeeID: 1,
		FullName:   "peerawas.k",
		Department: entity.Department{
			DepartmentID: 1,
			Name:         "Dev",
		},
		HourlyRate: hourRate,
	}

	timeSheet := entity.TimeSheet{
		HourWorked: decimal.NewFromFloat(8.0),
	}

	e := entity.PayStubs{
		PayStubID: 1,
		Employee:  emp,
		TaxAmount: decimal.NewFromFloat(0.07),
	}

	e.CalGrossPay(timeSheet.HourWorked)
	e.CalNetPay()

	expected := decimal.NewFromFloat(375.72)
	assert.True(t, expected.Equal(e.NetPay))
}
