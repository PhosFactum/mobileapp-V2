package patient_group

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *PatientGroupRepositoryImpl) GetPatientGroupsByCodeOrOrgTitle(search string, page, perPage int) ([]entities.PatientGroup, int64, error) {
	op := "repo.PatientGroup.SearchByCodeOrOrgTitle"
	var patientGroups []entities.PatientGroup
	var total int64

	baseQuery := r.db.
		Model(&entities.PatientGroup{}).
		Joins("JOIN organizations ON organizations.id = patient_groups.organization_id").
		Where("LOWER(patient_groups.code) LIKE LOWER(?) OR LOWER(organizations.title) LIKE LOWER(?)",
			"%"+search+"%", "%"+search+"%")

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

// GetPatientGroupsWithPatientsByOrganizationID получает группы пациентов с пациентами
func (r *PatientGroupRepositoryImpl) GetPatientGroupsWithPatientsByOrganizationID(orgID uint, page, perPage int) ([]entities.PatientGroup, int64, error) {
	op := "repo.PatientGroup.GetWithPatientsByOrganizationID"
	var patientGroups []entities.PatientGroup
	var total int64

	baseQuery := r.db.
		Model(&entities.PatientGroup{}).
		Where("organization_id = ?", orgID)

	// Получаем общее количество
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	// Получаем данные с пагинацией и предзагрузкой пациентов
	offset := (page - 1) * perPage
	err := baseQuery.
		Preload("Organization").
		Preload("Patient"). // Загружаем пациентов
		Offset(offset).
		Limit(perPage).
		Find(&patientGroups).
		Error

	if err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	return patientGroups, total, nil
}
