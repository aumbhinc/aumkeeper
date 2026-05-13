package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"aumkeeper/internal/domain"
	"aumkeeper/internal/repository"
)

type StaffService struct {
	repo repository.StaffRepository
}

func NewStaffService(
	repo repository.StaffRepository,
) *StaffService {
	return &StaffService{
		repo: repo,
	}
}

func (s *StaffService) CreateStaff(
	ctx context.Context,
	staff domain.Staff,
) (domain.Staff, error) {

	if strings.TrimSpace(staff.FirstName) == "" {
		return domain.Staff{}, fmt.Errorf("first name required")
	}

	if strings.TrimSpace(staff.LastName) == "" {
		return domain.Staff{}, fmt.Errorf("last name required")
	}

	if strings.TrimSpace(staff.Email) == "" {
		return domain.Staff{}, fmt.Errorf("email required")
	}

	if strings.TrimSpace(staff.Role) == "" {
		return domain.Staff{}, fmt.Errorf("role required")
	}

	staff.EmployeeCode = generateEmployeeCode()
	staff.Status = "active"

	return s.repo.Create(ctx, staff)
}

func generateEmployeeCode() string {
	return fmt.Sprintf(
		"AK-%d",
		time.Now().Unix(),
	)
}