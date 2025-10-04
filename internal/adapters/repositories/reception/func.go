package reception

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

// GetPatientReceptionsByPatientSpecializations - простой метод получения приемов
func (r *ReceptionRepositoryImpl) GetPatientReceptionsByPatientID(patientID uint) ([]entities.Reception, error) {
	op := "repo.Reception.GetPatientReceptionsByPatientID"

	var receptions []entities.Reception

	err := r.db.
		Where("patient_id = ?", patientID).
		Preload("Specialization").
		Order("created_at DESC").
		Find(&receptions).Error

	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return receptions, nil
}
