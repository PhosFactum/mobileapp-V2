package patient

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

// Транзакция прокидывается из юз кейса в репы
// Раздробить
// ДОПИЛИТЬ ВАЛИДАЦИЮ (создан ли с таким ФИО) + поля
// CreatePatient - создание пациента с кэшированием специальностей
func (r *PatientRepositoryImpl) CreatePatient(patientData *models.CreatePatientData, group_id uint) (*entities.Patient, error) {
	// 	op := "repo.Patient.CreatePatient"

	// 	var createdPatient *entities.Patient

	// 	err := r.db.Transaction(func(tx *gorm.DB) error {
	// 		// ✅ 1. Создаем контактную информацию
	// 		contactInfo := &entities.ContactInfo{
	// 			Phone:     patientData.ContactInfo.Phone,
	// 			Email:     patientData.ContactInfo.Email,
	// 			Address:   patientData.ContactInfo.Address,
	// 			CreatedAt: time.Now(),
	// 			UpdatedAt: time.Now(),
	// 		}
	// 		if err := tx.Create(contactInfo).Error; err != nil {
	// 			return fmt.Errorf("%s: failed to create contact info: %w", op, err)
	// 		}

	// 		// ✅ 2. Создаем персональную информацию
	// 		personalInfo := &entities.PersonalInfo{
	// 			DocNumber:      patientData.PersonalInfo.DocNumber,
	// 			DocSeries:      patientData.PersonalInfo.DocSeries,
	// 			SNILS:          patientData.PersonalInfo.SNILS,
	// 			OMS:            patientData.PersonalInfo.OMS,
	// 			DocumentTypeID: patientData.PersonalInfo.DocumentTypeID,
	// 			CreatedAt:      time.Now(),
	// 			UpdatedAt:      time.Now(),
	// 		}
	// 		if err := tx.Create(personalInfo).Error; err != nil {
	// 			return fmt.Errorf("%s: failed to create personal info: %w", op, err)
	// 		}

	// 		// ✅ 4. Создаем пустое направление на анализы
	// 		analysisOrder := &entities.AnalysisOrder{
	// 			OrderNumber: fmt.Sprintf("ORD-%06d", 0), // Временный номер
	// 			CreatedAt:   time.Now(),
	// 			UpdatedAt:   time.Now(),
	// 		}
	// 		if err := tx.Create(analysisOrder).Error; err != nil {
	// 			return fmt.Errorf("%s: failed to create analysis order: %w", op, err)
	// 		}

	// 		// Обновляем номер с правильным ID
	// 		analysisOrder.OrderNumber = fmt.Sprintf("ORD-%06d", analysisOrder.ID)
	// 		if err := tx.Save(analysisOrder).Error; err != nil {
	// 			return fmt.Errorf("%s: failed to update analysis order number: %w", op, err)
	// 		}

	// 		// ✅ 5. Создаем пациента
	// 		patient := &entities.Patient{
	// 			FullName:        patientData.FullName,
	// 			BirthDate:       patientData.BirthDate,
	// 			IsMale:          patientData.IsMale,
	// 			Position:        patientData.Position,
	// 			Division:        patientData.Division,
	// 			ExaminationType: patientData.ExaminationType,
	// 			ExaminationView: patientData.ExaminationView,
	// 			HarmPointID:     patientData.HarmPointID,
	// 			PatientGroupID:  group_id,
	// 			PersonalInfoID:  personalInfo.ID,
	// 			ContactInfoID:   contactInfo.ID,
	// 			AnalysisOrderID: analysisOrder.ID,
	// 			CreatedAt:       time.Now(),
	// 			UpdatedAt:       time.Now(),
	// 		}

	// 		if err := tx.Create(patient).Error; err != nil {
	// 			return fmt.Errorf("%s: failed to create patient %s: %w", op, patient.FullName, err)
	// 		}

	// 		// Обновляем AnalysisOrder с PatientID
	// 		analysisOrder.PatientID = patient.ID
	// 		if err := tx.Save(analysisOrder).Error; err != nil {
	// 			return fmt.Errorf("%s: failed to update analysis order with patient ID: %w", op, err)
	// 		}

	// 		// ✅ 6. КЭШИРУЕМ специальности пациента (определяем через HarmPoint)
	// 		specializations, err := r.getSpecializationsByHarmPoint(tx, patientData.HarmPointID)
	// 		if err != nil {
	// 			return fmt.Errorf("%s: failed to get specializations by harm point: %w", op, err)
	// 		}

	// 		// Связываем пациента со специализациями (кэшируем)
	// 		if err := tx.Model(patient).Association("Specializations").Append(&specializations); err != nil {
	// 			return fmt.Errorf("%s: failed to cache specializations: %w", op, err)
	// 		}

	// 		initialData := json.RawMessage(`{"values": {}, "schema": []}`)

	// 		// ✅ 8. Создаем пустые приемы (Reception) по каждой специализации
	// 		var receptions []entities.Reception
	// 		for _, spec := range specializations {
	// 			reception := entities.Reception{
	// 				PatientID:        patient.ID,
	// 				SpecializationID: spec.ID,
	// 				IsCompleted:      false,
	// 				CreatedAt:        time.Now(),
	// 				UpdatedAt:        time.Now(),
	// 				Data:             initialData,
	// 			}
	// 			receptions = append(receptions, reception)
	// 		}

	// 		if len(receptions) > 0 {
	// 			if err := tx.Create(&receptions).Error; err != nil {
	// 				return fmt.Errorf("%s: failed to create receptions: %w", op, err)
	// 			}
	// 		}

	// 		// ✅ 9. Создаем статистику пациента
	// 		statistics := &entities.PatientStatistics{
	// 			PatientID:              patient.ID,
	// 			TotalReceptions:        int64(len(specializations)),
	// 			CompletedReceptions:    0,
	// 			TotalAnalysisOrders:    0,
	// 			CompletedAnalysisItems: 0,
	// 			CreatedAt:              time.Now(),
	// 			UpdatedAt:              time.Now(),
	// 		}

	// 		if err := tx.Create(statistics).Error; err != nil {
	// 			if !strings.Contains(err.Error(), "duplicate") {
	// 				return fmt.Errorf("%s: failed to create statistics: %w", op, err)
	// 			}
	// 		}

	// 		// ✅ 10. Предзагружаем кэшированные специальности
	// 		tx.Preload("Specializations").First(patient, patient.ID)

	// 		// ✅ Сохраняем созданного пациента для возврата
	// 		createdPatient = patient

	// 		return nil
	// 	})

	// 	// ✅ Возвращаем пациента и ошибку
	// 	return createdPatient, err
	return nil, nil
}

// GetPatientsByGroup - получение пациентов по группе
func (r *PatientRepositoryImpl) GetPatientsByGroup(groupID uint) ([]entities.Patient, error) {
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

// getSpecializationsByHarmPoint - получение специальностей через HarmPoint
func (r *PatientRepositoryImpl) getSpecializationsByHarmPoint(tx *gorm.DB, harmPointID uint) ([]entities.Specialization, error) {
	op := "repo.Patient.getSpecializationsByHarmPoint"

	var specializations []entities.Specialization
	err := tx.Joins("JOIN harm_points_specializations hps ON hps.specialization_id = specializations.id").
		Where("hps.harm_point_id = ?", harmPointID).
		Find(&specializations).Error

	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return specializations, nil
}
