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

func (r *PatientGroupRepositoryImpl) GetPatientGroupsByDoctorID(doctorID uint, search string, page, perPage int) ([]entities.PatientGroup, int64, error) {
	op := "repo.PatientGroup.GetByDoctorID"
	var patientGroups []entities.PatientGroup
	var total int64

	// JOIN с промежуточной таблицей + организацией
	baseQuery := r.db.
		Model(&entities.PatientGroup{}).
		Joins("INNER JOIN doctor_patient_groups ON doctor_patient_groups.patient_group_id = patient_groups.id").
		Joins("INNER JOIN organizations ON organizations.id = patient_groups.organization_id").
		Where("doctor_patient_groups.doctor_id = ?", doctorID)

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		baseQuery = baseQuery.Where(
			"LOWER(patient_groups.code) LIKE ? OR LOWER(organizations.title) LIKE ?",
			searchPattern,
			searchPattern,
		)
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
