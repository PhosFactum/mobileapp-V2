package receptionHospital

import (
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

func (r *ReceptionHospitalRepositoryImpl) CreateReceptionHospital(reception entities.Reception) error {
	op := "repo.ReceptionHospital.CreateReceptionHospital"

	if err := r.db.Create(reception).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *ReceptionHospitalRepositoryImpl) UpdateReception(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.ReceptionHospital.UpdateReceptionHospital"

	var updatedReception entities.Reception
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedReception).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("hospital reception not found")
	}

	return updatedReception.ID, nil
}

func (r *ReceptionHospitalRepositoryImpl) DeleteReception(id uint) error {
	op := "repo.ReceptionHospital.DeleteReceptionHospital"

	result := r.db.Delete(&entities.Reception{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("hospital reception not found")
	}
	return nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionByID(id uint) (entities.Reception, error) {
	op := "repo.ReceptionHospital.GetReceptionHospitalByID"

	var reception entities.Reception
	if err := r.db.Preload("Doctor.Specialization").Preload("Patient").First(&reception, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Reception{}, errors.NewNotFoundError("hospital reception not found")
		}
		return entities.Reception{}, errors.NewDBError(op, err)
	}
	return reception, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionByPatientID(patientID uint) ([]entities.Reception, error) {
	op := "repo.ReceptionHospital.GetReceptionHospitalByDoctorID"

	var receptions []entities.Reception
	if err := r.db.
		Preload("Patient").
		Preload("Doctor.Specialization").
		Where("patient_id = ?", patientID).
		Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return receptions, nil
}

// // Обновлённая функция декодирования
// func decodeSpecializationData(data pgtype.JSONB, specialization string) (interface{}, error) {
// 	// Проверка наличия данных
// 	if data.Status != pgtype.Present || len(data.Bytes) == 0 {
// 		fmt.Print("BEDAAAA")
// 		return nil, nil
// 	}

// 	var result interface{}
// 	switch specialization {
// 	case "Невролог":
// 		result = new(entities.NeurologistData)
// 	case "Травматолог":
// 		result = new(entities.TraumatologistData)
// 	case "Психиатр":
// 		result = new(entities.PsychiatristData)
// 	case "Уролог":
// 		result = new(entities.UrologistData)
// 	case "Оториноларинголог":
// 		result = new(entities.OtolaryngologistData)
// 	case "Проктолог":
// 		result = new(entities.ProctologistData)
// 	case "Аллерголог":
// 		result = new(entities.AllergologistData)
// 	default:
// 		result = make(map[string]interface{})
// 	}

// 	// Декодирование с проверкой структуры
// 	decoder := json.NewDecoder(bytes.NewReader(data.Bytes))
// 	decoder.DisallowUnknownFields() // Для отлова несоответствий структур

// 	if err := decoder.Decode(&result); err != nil {
// 		log.Printf("Decoding error for %s: %v, data: %s", specialization, err, string(data.Bytes))
// 		return nil, fmt.Errorf("decoding error: %w", err)
// 	}

// 	return result, nil
// }
