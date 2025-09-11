package patientgroup

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *PatientGroupRepositoryImpl) GetByCodeOrOrgTitle(search string, page, perPage int) ([]entities.PatientGroup, int64, error) {
	op := "repo.PatientGroup.SearchByCodeOrOrgTitle"
	var patientGroups []entities.PatientGroup
	var total int64

	baseQuery := r.db.
		Model(&entities.PatientGroup{}).
		Joins("JOIN organizations ON organizations.id = patient_groups.organization_id").
		Where("patient_groups.code LIKE ? OR organizations.title LIKE ?",
			"%"+search+"%", "%"+search+"%")

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	offset := (page - 1) * perPage
	err := baseQuery.
		Offset(offset).
		Limit(perPage).
		Find(&patientGroups).
		Error

	if err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	return patientGroups, total, nil
}
