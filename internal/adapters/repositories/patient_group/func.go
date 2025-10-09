package patient_group

import (
	"strings"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

// GetPatientGroupsByOrganizationID возвращает группы пациентов конкретной организации
// с опциональным поиском по коду группы
func (r *PatientGroupRepositoryImpl) GetPatientGroupsByOrganizationID(orgID uint, search string, page, perPage int) ([]entities.PatientGroup, int64, error) {
	op := "repo.PatientGroup.GetByOrganizationID"
	var patientGroups []entities.PatientGroup
	var total int64

	baseQuery := r.db.
		Model(&entities.PatientGroup{}).
		Where("organization_id = ?", orgID)

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		baseQuery = baseQuery.Where("LOWER(code) LIKE ?", searchPattern)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	offset := (page - 1) * perPage
	err := baseQuery.
		Preload("Organization").
		Offset(offset).
		Limit(perPage).
		Find(&patientGroups).
		Error

	if err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	return patientGroups, total, nil
}
