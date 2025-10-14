// repositories/consent_signatures/func.go
package consent_signatures

import (
	"log"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

// SaveSignature сохраняет подпись пациента
func (r *ConsentSignatureRepositoryImpl) SaveSignature(patientID uint, signature []byte) error {
	op := "repo.ConsentSignature.SaveSignature"

	// Пытаемся найти существующую запись
	var existing entities.ConsentSignature
	err := r.db.Where("patient_id = ?", patientID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Создаём новую запись
		newSignature := entities.ConsentSignature{
			PatientID: patientID,
			Signature: signature,
		}
		if createErr := r.db.Create(&newSignature).Error; createErr != nil {
			log.Printf("[%s] Failed to create signature for patient %d: %v", op, patientID, createErr)
			return errors.NewDBError(op, createErr)
		}
	} else if err != nil {
		// Другая ошибка БД
		log.Printf("[%s] DB error: %v", op, err)
		return errors.NewDBError(op, err)
	} else {
		// Обновляем существующую запись
		existing.Signature = signature
		if updateErr := r.db.Save(&existing).Error; updateErr != nil {
			log.Printf("[%s] Failed to update signature for patient %d: %v", op, patientID, updateErr)
			return errors.NewDBError(op, updateErr)
		}
	}

	log.Printf("[%s] Signature saved for patient %d", op, patientID)
	return nil
}

// GetSignature возвращает подпись по ID пациента
func (r *ConsentSignatureRepositoryImpl) GetSignature(patientID uint) ([]byte, error) {
	op := "repo.ConsentSignature.GetSignature"

	var signature entities.ConsentSignature
	if err := r.db.Select("signature").
		Where("patient_id = ?", patientID).
		First(&signature).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			log.Printf("[%s] Signature not found for patient %d", op, patientID)
			return nil, errors.NewNotFoundError(op)
		}

		log.Printf("[%s] DB error: %v", op, err)
		return nil, errors.NewDBError(op, err)
	}

	log.Printf("[%s] Signature fetched for patient %d", op, patientID)
	return signature.Signature, nil
}

// GetByPatientID возвращает полную сущность
func (r *ConsentSignatureRepositoryImpl) GetByPatientID(patientID uint) (*entities.ConsentSignature, error) {
	op := "repo.ConsentSignature.GetByPatientID"

	var signature entities.ConsentSignature
	if err := r.db.Where("patient_id = ?", patientID).First(&signature).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Не найдено — не ошибка
		}
		log.Printf("[%s] DB error: %v", op, err)
		return nil, errors.NewDBError(op, err)
	}

	return &signature, nil
}
