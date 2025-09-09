package organization

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *OrganizationRepositoryImpl) GetAllOrganizations(page, perPage int) ([]entities.Organization, int64, error) {
	op := "repo.Organization.GetAllOrganizations"
	var organizations []entities.Organization
	var total int64

	baseQuery := r.db.
		Model(&entities.Organization{})

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	offset := (page - 1) * perPage
	err := baseQuery.
		Preload("Doctor").
		Offset(offset).
		Limit(perPage).
		Find(&organizations).
		Error

	if err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	return organizations, total, nil
}

func (r *OrganizationRepositoryImpl) GetByTitleOrCodeOrganizations(search string, page, perPage int) ([]entities.Organization, int64, error) {
	op := "repo.Organization.GetByTitleOrCodeOrganizations"
	var organizations []entities.Organization
	var total int64

	baseQuery := r.db.
		Model(&entities.Organization{})

	// Добавляем поиск по началу подстроки Title или Code
	if search != "" {
		baseQuery = baseQuery.Where(
			"LOWER(title) LIKE LOWER(?) OR LOWER(code) LIKE LOWER(?)",
			search+"%",
			search+"%",
		)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	offset := (page - 1) * perPage
	err := baseQuery.
		Preload("Doctor").
		Offset(offset).
		Limit(perPage).
		Find(&organizations).
		Error

	if err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	return organizations, total, nil
}
