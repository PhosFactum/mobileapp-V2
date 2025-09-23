package patient

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *PatientRepositoryImpl) GetPatientByIDWithTx(tx *gorm.DB, id uint) (*entities.Patient, error) {
	op := "repo.Patient.UpdatePatientWithTx"
	var patient entities.Patient
	if err := tx.First(&patient, id).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return &patient, nil
}

func (r *PatientRepositoryImpl) UpdatePatientWithTx(tx *gorm.DB, id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.Patient.UpdatePatientWithTx"

	var updatedPatient entities.Patient
	result := tx.
		Clauses(clause.Returning{}).
		Model(&updatedPatient).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("patient not found")
	}

	return updatedPatient.ID, nil
}

func (r *PatientRepositoryImpl) UpdatePatient(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.Patient.UpdatePatient"

	var updatedPatient entities.Patient
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedPatient).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewDBError(op, result.Error)
	}

	return updatedPatient.ID, nil
}

func (r *PatientRepositoryImpl) DeletePatient(id uint) error {
	op := "repo.Patient.ForceDeletePatient"

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Удаляем самые "глубокие" связи
	if err := tx.Exec("DELETE FROM reception_smp_med_services WHERE reception_smp_id IN (SELECT id FROM reception_smps WHERE patient_id = ?)", id).Error; err != nil {
		tx.Rollback()
		return errors.NewDBError(op, fmt.Errorf("failed to delete from reception_smp_med_services: %w", err))
	}

	// 2. Удаляем зависимости 1 уровня
	if err := tx.Exec("DELETE FROM reception_smps WHERE patient_id = ?", id).Error; err != nil {
		tx.Rollback()
		return errors.NewDBError(op, fmt.Errorf("failed to delete from reception_smps: %w", err))
	}

	if err := tx.Exec("DELETE FROM reception_hospitals WHERE patient_id = ?", id).Error; err != nil {
		tx.Rollback()
		return errors.NewDBError(op, fmt.Errorf("failed to delete from reception_hospitals: %w", err))
	}

	if err := tx.Exec("DELETE FROM patient_allergy WHERE patient_id = ?", id).Error; err != nil {
		tx.Rollback()
		return errors.NewDBError(op, fmt.Errorf("failed to delete from patient_allergy: %w", err))
	}

	// 3. Удаляем самого пациента
	result := tx.Delete(&entities.Patient{}, id)
	if result.Error != nil {
		tx.Rollback()
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.NewDBError(op, errors.NewNotFoundError("patient not found"))
	}

	// 4. Коммитим транзакцию
	if err := tx.Commit().Error; err != nil {
		return errors.NewDBError(op, fmt.Errorf("failed to commit transaction: %w", err))
	}

	return nil
}

func (r *PatientRepositoryImpl) GetPatientByID(id uint) (entities.Patient, error) {
	op := "repo.Patient.GetPatientByID"

	var patient entities.Patient
	err := r.db.
		Preload("PersonalInfo").
		Preload("ContactInfo").
		Preload("Allergy").
		First(&patient, id).Error

	if err == gorm.ErrRecordNotFound {
		return entities.Patient{}, errors.NewNotFoundError("patient not found")
	}

	if err != nil {
		return entities.Patient{}, errors.NewDBError(op, err)
	}

	return patient, nil
}

func (r *PatientRepositoryImpl) GetPatientsByFullName(name string) ([]entities.Patient, error) {
	op := "repo.Patient.GetPatientsByFullName"

	var patients []entities.Patient
	if err := r.db.
		Where("full_name ILIKE ?", "%"+name+"%").
		Preload("PersonalInfo").
		Preload("ContactInfo").
		Find(&patients).Error; err != nil {

		return nil, errors.NewDBError(op, err)
	}

	if len(patients) == 0 {
		return nil, errors.NewDBError(op, errors.NewNotFoundError("no patients found"))
	}

	return patients, nil
}

func (r *PatientRepositoryImpl) GetAllPatients(page, count int, queryFilter string, queryOrder string, parameters []interface{}) ([]entities.Patient, int64, error) {
	op := "repo.Patient.GetAllPatients"

	// Создаем базовый запрос
	query := r.db.Model(&entities.Patient{})

	// Применяем фильтрацию
	if queryFilter != "" {
		query = query.Where(queryFilter, parameters...)
	}

	// Применяем сортировку
	if queryOrder != "" {
		query = query.Order(queryOrder)
	}

	// Подсчитываем общее количество записей
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	// Применяем пагинацию
	if page > 0 && count > 0 {
		offset := (page - 1) * count
		query = query.Offset(offset).Limit(count)
	}

	// Получаем записи
	var patients []entities.Patient
	result := query.Find(&patients)
	if result.Error != nil {
		return nil, 0, errors.NewDBError(op, result.Error)
	}

	return patients, totalRecords, nil
}

// GetPatientsByGroup - получение пациентов по группе
func (r *PatientRepositoryImpl) GetPatientsByGroup(page, perPage int, group_id uint) ([]entities.Patient, int64, error) {
	op := "repo.Patient.GetPatientsByGroup"

	// Подсчитываем общее количество пациентов в группе
	var totalRecords int64
	countQuery := r.db.Model(&entities.Patient{}).
		Where("patient_group_id = ?", group_id)

	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	// Основной запрос с предзагрузкой
	query := r.db.Model(&entities.Patient{}).
		Preload("ExaminationType").
		Preload("ExaminationView").
		Preload("HarmPoint").
		Preload("PersonalInfo.DocumentType").
		Preload("ContactInfo").
		Preload("Flg").
		Preload("Organization").
		Preload("AnalysisOrder.OrderItems.Analysis").
		Preload("Vaccines").
		Preload("Receptions.Specialization").
		Preload("Specializations").
		Preload("Statistics").
		Where("patient_group_id = ?", group_id).
		Order("full_name")

	// Пагинация
	if page > 0 && perPage > 0 {
		offset := (page - 1) * perPage
		query = query.Offset(offset).Limit(perPage)
	}

	var patients []entities.Patient
	result := query.Find(&patients)
	if result.Error != nil {
		return nil, 0, errors.NewDBError(op, result.Error)
	}

	return patients, totalRecords, nil
}

// ДОПИЛИТЬ ВАЛИДАЦИЮ (создан ли с таким ФИО) + поля
// CreatePatient - создание пациента с кэшированием специальностей
func (r *PatientRepositoryImpl) CreatePatient(patientData *models.CreatePatientData, group_id uint) (*entities.Patient, error) {
	op := "repo.Patient.CreatePatient"

	var createdPatient *entities.Patient

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// ✅ 1. Создаем контактную информацию
		contactInfo := &entities.ContactInfo{
			Phone:     patientData.ContactInfo.Phone,
			Email:     patientData.ContactInfo.Email,
			Address:   patientData.ContactInfo.Address,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := tx.Create(contactInfo).Error; err != nil {
			return fmt.Errorf("%s: failed to create contact info: %w", op, err)
		}

		// ✅ 2. Создаем персональную информацию
		personalInfo := &entities.PersonalInfo{
			DocNumber:      patientData.PersonalInfo.DocNumber,
			DocSeries:      patientData.PersonalInfo.DocSeries,
			SNILS:          patientData.PersonalInfo.SNILS,
			OMS:            patientData.PersonalInfo.OMS,
			DocumentTypeID: patientData.PersonalInfo.DocumentTypeID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := tx.Create(personalInfo).Error; err != nil {
			return fmt.Errorf("%s: failed to create personal info: %w", op, err)
		}

		// ✅ 4. Создаем пустое направление на анализы
		analysisOrder := &entities.AnalysisOrder{
			OrderNumber: fmt.Sprintf("ORD-%06d", 0), // Временный номер
			TotalAmount: 0,                          // Пустое направление
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := tx.Create(analysisOrder).Error; err != nil {
			return fmt.Errorf("%s: failed to create analysis order: %w", op, err)
		}

		// Обновляем номер с правильным ID
		analysisOrder.OrderNumber = fmt.Sprintf("ORD-%06d", analysisOrder.ID)
		if err := tx.Save(analysisOrder).Error; err != nil {
			return fmt.Errorf("%s: failed to update analysis order number: %w", op, err)
		}

		// ✅ 5. Создаем пациента
		patient := &entities.Patient{
			FullName:          patientData.FullName,
			BirthDate:         patientData.BirthDate,
			IsMale:            patientData.IsMale,
			Position:          patientData.Position,
			Division:          patientData.Division,
			ExaminationTypeID: patientData.ExaminationTypeID,
			ExaminationViewID: patientData.ExaminationViewID,
			HarmPointID:       patientData.HarmPointID,
			PatientGroupID:    group_id,
			PersonalInfoID:    personalInfo.ID,
			ContactInfoID:     contactInfo.ID,
			AnalysisOrderID:   analysisOrder.ID,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := tx.Create(patient).Error; err != nil {
			return fmt.Errorf("%s: failed to create patient %s: %w", op, patient.FullName, err)
		}

		// Обновляем AnalysisOrder с PatientID
		analysisOrder.PatientID = patient.ID
		if err := tx.Save(analysisOrder).Error; err != nil {
			return fmt.Errorf("%s: failed to update analysis order with patient ID: %w", op, err)
		}

		// ✅ 6. КЭШИРУЕМ специальности пациента (определяем через HarmPoint)
		specializations, err := r.getSpecializationsByHarmPoint(tx, patientData.HarmPointID)
		if err != nil {
			return fmt.Errorf("%s: failed to get specializations by harm point: %w", op, err)
		}

		// Связываем пациента со специализациями (кэшируем)
		if err := tx.Model(patient).Association("Specializations").Append(&specializations); err != nil {
			return fmt.Errorf("%s: failed to cache specializations: %w", op, err)
		}

		// ✅ 8. Создаем пустые приемы (Reception) по каждой специализации
		var receptions []entities.Reception
		for _, spec := range specializations {
			reception := entities.Reception{
				PatientID:        patient.ID,
				SpecializationID: spec.ID,
				IsCompleted:      false,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
				SpecializationData: pgtype.JSONB{
					Status: pgtype.Null,
				},
				CustomFieldsSchema: nil,
			}
			receptions = append(receptions, reception)
		}

		if len(receptions) > 0 {
			if err := tx.Create(&receptions).Error; err != nil {
				return fmt.Errorf("%s: failed to create receptions: %w", op, err)
			}
		}

		// ✅ 9. Создаем статистику пациента
		statistics := &entities.PatientStatistics{
			PatientID:              patient.ID,
			TotalReceptions:        int64(len(specializations)),
			CompletedReceptions:    0,
			TotalAnalysisOrders:    0,
			CompletedAnalysisItems: 0,
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
		}

		if err := tx.Create(statistics).Error; err != nil {
			if !strings.Contains(err.Error(), "duplicate") {
				return fmt.Errorf("%s: failed to create statistics: %w", op, err)
			}
		}

		// ✅ 10. Предзагружаем кэшированные специальности
		tx.Preload("Specializations").First(patient, patient.ID)

		// ✅ Сохраняем созданного пациента для возврата
		createdPatient = patient

		return nil
	})

	// ✅ Возвращаем пациента и ошибку
	return createdPatient, err
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
