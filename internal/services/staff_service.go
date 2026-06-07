package services

import (
	"context"
	"fmt"
	"time"

	"aumkeeper/internal/domain"
	"aumkeeper/internal/repository"
)

type StaffService struct {
	repo repository.StaffRepository
}

func NewStaffService(repo repository.StaffRepository) *StaffService {
	return &StaffService{repo: repo}
}

func (s *StaffService) CreateStaff(ctx context.Context, staff domain.Staff) (domain.Staff, error) {

	// ==============================
	// 1. REQUIRED SYSTEM FIELDS
	// ==============================

	if staff.FirstName == "" || staff.LastName == "" {
		return domain.Staff{}, fmt.Errorf("first name and last name required")
	}

	if staff.PhoneNumber == "" {
		return domain.Staff{}, fmt.Errorf("phone number required")
	}

	// ==============================
	// 2. SYSTEM-GENERATED FIELDS
	// ==============================

	if staff.EmployeeCode == "" {
		staff.EmployeeCode = generateEmployeeCode()
	}

	if staff.Status == "" {
		staff.Status = "active"
	}

	if staff.TaxFileStatus == "" {
		staff.TaxFileStatus = "not_filed"
	}

	// ==============================
	// 3. SAFE DEFAULTS
	// ==============================

	if staff.PaymentFrequency == "" {
		staff.PaymentFrequency = "monthly"
	}

	// ==============================
	// 4. CREATE TIMESTAMP SAFETY (optional override safety)
	// ==============================

	if staff.CreatedAt.IsZero() {
		staff.CreatedAt = time.Now()
	}

	// ==============================
	// 5. PERSIST
	// ==============================

	return s.repo.Create(ctx, staff)
}

// Simple but safe generator (replace later with SKU service)
func generateEmployeeCode() string {
	return fmt.Sprintf("EMP-%d", time.Now().UnixNano())
}