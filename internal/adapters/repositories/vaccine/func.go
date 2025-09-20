package vaccine

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

// CreateVaccine - создание вакцины
func (r *VaccineRepository) CreateVaccine(vaccine *entities.Vaccine) error {
	op := "repo.Vaccine.CreateVaccine"

	if err := r.db.Create(vaccine).Error; err != nil {
		return errors.NewDBError(op, err)
	}

	return nil
}

// UpdateVaccine - обновление вакцины
func (r *VaccineRepository) UpdateVaccine(vaccine *entities.Vaccine) error {
	op := "repo.Vaccine.UpdateVaccine"

	if err := r.db.Save(vaccine).Error; err != nil {
		return errors.NewDBError(op, err)
	}

	return nil
}

// DeleteVaccine - удаление вакцины
func (r *VaccineRepository) DeleteVaccine(patientID, vaccineID uint) error {
	op := "repo.Vaccine.DeleteVaccine"

	result := r.db.Where("patient_id = ? AND id = ?", patientID, vaccineID).
		Delete(&entities.Vaccine{})

	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("no such vaccine to delete")
	}

	return nil
}

// GetPatientVaccines - получение всех вакцин пациента
func (r *VaccineRepository) GetPatientVaccines(patientID uint) ([]entities.Vaccine, error) {
	op := "repo.Vaccine.GetPatientVaccines"

	var vaccines []entities.Vaccine
	err := r.db.Where("patient_id = ?", patientID).
		Order("date DESC").
		Find(&vaccines).Error

	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return vaccines, nil
}
