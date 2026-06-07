package domain

import "time"

// Staff is the core workforce entity in AumKeeper ERP.
// It represents a full employee record used across:
// - onboarding
// - payroll
// - scheduling
// - execution tracking
type Staff struct {
	ID           int
	EmployeeCode string

	// Identity
	FirstName  string
	MiddleName string
	LastName   string

	// Job
	Role string

	// Contact
	Email       string
	PhoneNumber string

	// Address
	Street string
	City   string
	State  string
	ZipCode string

	// Compliance / Payroll
	SSN              string
	TaxFileStatus    string
	DependentClaims  int
	Wage             float64
	PaymentFrequency string

	// Internal
	Comments string
	Status   string

	// System timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}