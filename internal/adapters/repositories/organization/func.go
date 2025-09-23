package organization

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *OrganizationRepositoryImpl) GetAllOrganizations(doctorID uint, page, perPage int) ([]entities.Organization, int64, error) {
	op := "repo.Organization.GetAllOrganizations"
	var organizations []entities.Organization
	var total int64

	baseQuery := r.db.
		Model(&entities.Organization{})

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	offset := (page - 1) * perPage

	err := r.db.
		Table("organizations").
		Joins("JOIN doctor_organizations ON doctor_organizations.organization_id = organizations.id").
		Where("doctor_organizations.doctor_id = ?", doctorID).
		Preload("Manager").
		Offset(offset).
		Limit(perPage).
		Find(&organizations).Error

	if err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	return organizations, total, nil
}
