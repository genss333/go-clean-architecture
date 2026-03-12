package entity

import "github.com/shopspring/decimal"

type PayStubs struct {
	PayStubID int
	Employee  Employee
	GrossPay  decimal.Decimal
	TaxAmount decimal.Decimal
	NetPay    decimal.Decimal
}

func (p *PayStubs) CalGrossPay(hourWorked decimal.Decimal) {
	houryRate := p.Employee.HourlyRate.Amount
	p.GrossPay = houryRate.Mul(hourWorked)
}

func (p *PayStubs) CalNetPay() {
	tax := p.GrossPay.Mul(p.TaxAmount)
	p.NetPay = p.GrossPay.Sub(tax)
}
