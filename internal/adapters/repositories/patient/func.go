package patient

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

// GetPatientsByGroup - получение пациентов по группе
func (r *PatientRepositoryImpl) GetPatientsByGroup(groupID uint) ([]entities.Patient, *errors.AppError) {
	op := "repo.Patient.GetPatientsByGroup"

	query := r.db.
		Preload("HarmPoint").
		Preload("PersonalInfo").
		Preload("ContactInfo").
		Preload("Flg").
		Preload("AnalysisOrder.OrderItems.Analysis").
		Preload("Vaccines").
		Preload("VaccineRefusals").
		Preload("VaccineWithdrawals").
		Preload("Titers").
		Preload("Receptions.Specialization").
		Preload("Receptions.Template").
		Preload("Specializations").
		Preload("Statistics").
		Where("patient_group_id = ?", groupID).
		Order("full_name")

	var patients []entities.Patient
	if err := query.Find(&patients).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return patients, nil
}

// CreateContactInfo создаёт контактную информацию
func (r *PatientRepositoryImpl) CreateContactInfo(tx *gorm.DB, contactInfo *entities.ContactInfo) *errors.AppError {
	op := "repo.Patient.CreateContactInfo"
	if err := tx.Create(contactInfo).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreatePersonalInfo создаёт персональную информацию
func (r *PatientRepositoryImpl) CreatePersonalInfo(tx *gorm.DB, personalInfo *entities.PersonalInfo) *errors.AppError {
	op := "repo.Patient.CreatePersonalInfo"
	if err := tx.Create(personalInfo).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreateAnalysisOrder создаёт направление на анализы
func (r *PatientRepositoryImpl) CreateAnalysisOrder(tx *gorm.DB, order *entities.AnalysisOrder) *errors.AppError {
	op := "repo.Patient.CreateAnalysisOrder"
	if err := tx.Create(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// UpdateAnalysisOrder обновляет направление
func (r *PatientRepositoryImpl) UpdateAnalysisOrder(tx *gorm.DB, order *entities.AnalysisOrder) *errors.AppError {
	op := "repo.Patient.UpdateAnalysisOrder"
	if err := tx.Save(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreatePatient создаёт пациента
func (r *PatientRepositoryImpl) CreatePatient(tx *gorm.DB, patient *entities.Patient) *errors.AppError {
	op := "repo.Patient.CreatePatient"
	if err := tx.Create(patient).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CacheSpecializations связывает пациента со специализациями
func (r *PatientRepositoryImpl) CacheSpecializations(tx *gorm.DB, patient *entities.Patient, specializations []entities.Specialization) *errors.AppError {
	op := "repo.Patient.CacheSpecializations"
	if err := tx.Model(patient).Association("Specializations").Append(&specializations); err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreateReceptions создаёт приёмы
func (r *PatientRepositoryImpl) CreateReceptions(tx *gorm.DB, receptions []entities.Reception) *errors.AppError {
	if len(receptions) == 0 {
		return nil
	}
	op := "repo.Patient.CreateReceptions"
	if err := tx.Create(&receptions).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreatePatientStatistics создаёт статистику пациента
func (r *PatientRepositoryImpl) CreatePatientStatistics(tx *gorm.DB, stats *entities.PatientStatistics) *errors.AppError {
	op := "repo.Patient.CreatePatientStatistics"
	if err := tx.Create(stats).Error; err != nil {
		return errors.NewDBError(op, err) // ← стандартная обработка, как просили
	}
	return nil
}

// GetReceptionTemplatesByHarmPointID возвращает шаблоны заключений по HarmPointID
func (r *PatientRepositoryImpl) GetReceptionTemplatesByHarmPointID(tx *gorm.DB, harmPointID uint) ([]entities.ReceptionTemplate, *errors.AppError) {
	op := "repo.Patient.GetReceptionTemplatesByHarmPointID"

	// Загружаем HarmPoint с предзагрузкой шаблонов
	var harmPoint entities.HarmPoint
	if err := tx.Preload("ReceptionTemplates").First(&harmPoint, harmPointID).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return harmPoint.ReceptionTemplates, nil
}

// PreloadPatientWithSpecializations загружает пациента со специализациями
func (r *PatientRepositoryImpl) PreloadPatientWithSpecializations(tx *gorm.DB, patientID uint) (*entities.Patient, *errors.AppError) {
	op := "repo.Patient.PreloadPatientWithSpecializations"
	var patient entities.Patient
	if err := tx.Preload("Specializations").First(&patient, patientID).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return &patient, nil
}

// func (r *PatientRepositoryImpl) GetReceptionTemplatesByHarmPointFromModel(
// 	tx *gorm.DB,
// 	harmPoint *entities.HarmPoint,
// ) ([]entities.ReceptionTemplate, error) {
// 	op := "repo.Patient.GetReceptionTemplatesByHarmPointFromModel"
// 	var templates []entities.ReceptionTemplate
// 	if err := tx.Model(harmPoint).Association("ReceptionTemplates").Find(&templates); err != nil {
// 		return nil, errors.NewDBError(op, err)
// 	}
// 	return templates, nil
// }
