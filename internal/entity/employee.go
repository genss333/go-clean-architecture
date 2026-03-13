package entity

type Employee struct {
	EmployeeID int
	FullName   string

	Department Department
	HourlyRate HourlyRate
}
