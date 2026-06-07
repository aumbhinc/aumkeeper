package repository

import (
	"context"
	"fmt"

	"aumkeeper/internal/domain"
)

type StaffRepository interface {
	Create(
		ctx context.Context,
		staff domain.Staff,
	) (domain.Staff, error)
}

type staffRepository struct {
	repos *Repositories
}

func NewStaffRepository(r *Repositories) StaffRepository {
	return &staffRepository{
		repos: r,
	}
}

func (r *staffRepository) Create(
	ctx context.Context,
	staff domain.Staff,
) (domain.Staff, error) {

	query := `
		INSERT INTO staff (
			employee_code,
			first_name,
			middle_name,
			last_name,
			role,
			email,
			phone_number,
			street,
			city,
			state,
			zip_code,
			ssn,
			tax_file_status,
			dependent_claims,
			wage,
			payment_frequency,
			comments,
			status
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,
			$10,$11,$12,$13,$14,$15,$16,
			$17,$18
		)
		RETURNING id, created_at, updated_at;
	`

	err := r.repos.DB.QueryRow(
		ctx,
		query,

		staff.EmployeeCode,

		staff.FirstName,
		staff.MiddleName,
		staff.LastName,

		staff.Role,

		staff.Email,
		staff.PhoneNumber,

		staff.Street,
		staff.City,
		staff.State,
		staff.ZipCode,

		staff.SSN,
		staff.TaxFileStatus,
		staff.DependentClaims,

		staff.Wage,
		staff.PaymentFrequency,

		staff.Comments,

		staff.Status,
	).Scan(
		&staff.ID,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)

	if err != nil {
		return domain.Staff{}, fmt.Errorf("create staff: %w", err)
	}

	return staff, nil
}