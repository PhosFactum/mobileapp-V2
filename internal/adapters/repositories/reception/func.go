package receptionHospital

// // GetDoctorReceptions - получение приемов доступных доктору (по его специализациям)
// func (r *ReceptionRepositoryImpl) GetDoctorReceptions(doctorID uint) ([]entities.Reception, error) {
//     op := "repo.Reception.GetDoctorReceptions"

//     var receptions []entities.Reception
//     err := r.db.Joins("JOIN doctor_specializations ds ON ds.specialization_id = receptions.specialization_id").
//         Where("ds.doctor_id = ?", doctorID).
//         Order("receptions.created_at DESC").
//         Find(&receptions).Error

//     if err != nil {
//         return nil, errors.NewDBError(op, err)
//     }

//     return receptions, nil
// }

// func (r *ReceptionHospitalRepositoryImpl) CreateReceptionHospital(reception entities.Reception) error {
// 	op := "repo.ReceptionHospital.CreateReceptionHospital"

// 	if err := r.db.Create(reception).Error; err != nil {
// 		return errors.NewDBError(op, err)
// 	}
// 	return nil
// }

// func (r *ReceptionHospitalRepositoryImpl) UpdateReception(id uint, updateMap map[string]interface{}) (uint, error) {
// 	op := "repo.ReceptionHospital.UpdateReceptionHospital"

// 	var updatedReception entities.Reception
// 	result := r.db.
// 		Clauses(clause.Returning{}).
// 		Model(&updatedReception).
// 		Where("id = ?", id).
// 		Updates(updateMap)

// 	if result.Error != nil {
// 		return 0, errors.NewDBError(op, result.Error)
// 	}
// 	if result.RowsAffected == 0 {
// 		return 0, errors.NewNotFoundError("hospital reception not found")
// 	}

// 	return updatedReception.ID, nil
// }

// func (r *ReceptionHospitalRepositoryImpl) DeleteReception(id uint) error {
// 	op := "repo.ReceptionHospital.DeleteReceptionHospital"

// 	result := r.db.Delete(&entities.Reception{}, id)
// 	if result.Error != nil {
// 		return errors.NewDBError(op, result.Error)
// 	}
// 	if result.RowsAffected == 0 {
// 		return errors.NewNotFoundError("hospital reception not found")
// 	}
// 	return nil
// }

// func (r *ReceptionHospitalRepositoryImpl) GetReceptionByID(id uint) (entities.Reception, error) {
// 	op := "repo.ReceptionHospital.GetReceptionHospitalByID"

// 	var reception entities.Reception
// 	if err := r.db.Preload("Doctor.Specialization").Preload("Patient").First(&reception, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return entities.Reception{}, errors.NewNotFoundError("hospital reception not found")
// 		}
// 		return entities.Reception{}, errors.NewDBError(op, err)
// 	}
// 	return reception, nil
// }

// func (r *ReceptionHospitalRepositoryImpl) GetReceptionByPatientID(patientID uint) ([]entities.Reception, error) {
// 	op := "repo.ReceptionHospital.GetReceptionHospitalByDoctorID"

// 	var receptions []entities.Reception
// 	if err := r.db.
// 		Preload("Patient").
// 		Preload("Doctor.Specialization").
// 		Where("patient_id = ?", patientID).
// 		Find(&receptions).Error; err != nil {
// 		return nil, errors.NewDBError(op, err)
// 	}

// 	return receptions, nil
// }

// // // Обновлённая функция декодирования
// // func decodeSpecializationData(data pgtype.JSONB, specialization string) (interface{}, error) {
// // 	// Проверка наличия данных
// // 	if data.Status != pgtype.Present || len(data.Bytes) == 0 {
// // 		fmt.Print("BEDAAAA")
// // 		return nil, nil
// // 	}

// // 	var result interface{}
// // 	switch specialization {
// // 	case "Невролог":
// // 		result = new(entities.NeurologistData)
// // 	case "Травматолог":
// // 		result = new(entities.TraumatologistData)
// // 	case "Психиатр":
// // 		result = new(entities.PsychiatristData)
// // 	case "Уролог":
// // 		result = new(entities.UrologistData)
// // 	case "Оториноларинголог":
// // 		result = new(entities.OtolaryngologistData)
// // 	case "Проктолог":
// // 		result = new(entities.ProctologistData)
// // 	case "Аллерголог":
// // 		result = new(entities.AllergologistData)
// // 	default:
// // 		result = make(map[string]interface{})
// // 	}

// // 	// Декодирование с проверкой структуры
// // 	decoder := json.NewDecoder(bytes.NewReader(data.Bytes))
// // 	decoder.DisallowUnknownFields() // Для отлова несоответствий структур

// // 	if err := decoder.Decode(&result); err != nil {
// // 		log.Printf("Decoding error for %s: %v, data: %s", specialization, err, string(data.Bytes))
// // 		return nil, fmt.Errorf("decoding error: %w", err)
// // 	}

// // 	return result, nil
// // }

// // repository/reception_repository.go

// // GetPatientReceptions - получение всех приемов пациента
// func (r *ReceptionRepository) GetPatientReceptions(patientID uint) ([]entities.Reception, error) {
//     op := "repo.Reception.GetPatientReceptions"

//     var receptions []entities.Reception
//     err := r.db.Preload("Specialization").
//         Where("patient_id = ?", patientID).
//         Order("created_at DESC").
//         Find(&receptions).Error

//     if err != nil {
//         return nil, errors.NewDBError(op, err)
//     }

//     return receptions, nil
// }

// // GetPatientReceptionByID - получение конкретного приема пациента
// func (r *ReceptionRepository) GetPatientReceptionByID(patientID, receptionID uint) (*entities.Reception, error) {
//     op := "repo.Reception.GetPatientReceptionByID"

//     var reception entities.Reception
//     err := r.db.Preload("Specialization").
//         Where("patient_id = ? AND id = ?", patientID, receptionID).
//         First(&reception).Error

//     if err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             return nil, fmt.Errorf("%s: reception not found", op)
//         }
//         return nil, errors.NewDBError(op, err)
//     }

//     return &reception, nil
// }

// // GetPatientReceptionBySpecialization - получение приема по специализации
// func (r *ReceptionRepository) GetPatientReceptionBySpecialization(patientID, specializationID uint) (*entities.Reception, error) {
//     op := "repo.Reception.GetPatientReceptionBySpecialization"

//     var reception entities.Reception
//     err := r.db.Preload("Specialization").
//         Where("patient_id = ? AND specialization_id = ?", patientID, specializationID).
//         First(&reception).Error

//     if err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             return nil, fmt.Errorf("%s: reception not found for this specialization", op)
//         }
//         return nil, errors.NewDBError(op, err)
//     }

//     return &reception, nil
// }

// // CreateReception - создание приема
// func (r *ReceptionRepository) CreateReception(reception *entities.Reception) error {
//     op := "repo.Reception.CreateReception"

//     // Проверяем уникальное ограничение
//     var existing entities.Reception
//     err := r.db.Where("patient_id = ? AND specialization_id = ?",
//         reception.PatientID, reception.SpecializationID).
//         First(&existing).Error

//     if err == nil {
//         return fmt.Errorf("%s: reception for this specialization already exists", op)
//     }

//     if !errors.Is(err, gorm.ErrRecordNotFound) {
//         return errors.NewDBError(op, err)
//     }

//     if err := r.db.Create(reception).Error; err != nil {
//         return errors.NewDBError(op, err)
//     }

//     return nil
// }

// // UpdateReception - обновление приема
// func (r *ReceptionRepository) UpdateReception(reception *entities.Reception) error {
//     op := "repo.Reception.UpdateReception"

//     if err := r.db.Save(reception).Error; err != nil {
//         return errors.NewDBError(op, err)
//     }

//     return nil
// }

// // DeleteReception - удаление приема
// func (r *ReceptionRepository) DeleteReception(patientID, receptionID uint) error {
//     op := "repo.Reception.DeleteReception"

//     result := r.db.Where("patient_id = ? AND id = ?", patientID, receptionID).
//         Delete(&entities.Reception{})

//     if result.Error != nil {
//         return errors.NewDBError(op, result.Error)
//     }

//     if result.RowsAffected == 0 {
//         return fmt.Errorf("%s: reception not found", op)
//     }

//     return nil
// }

// // GetPatientReceptionsWithPagination - получение приемов с пагинацией
// func (r *ReceptionRepository) GetPatientReceptionsWithPagination(
//     patientID uint,
//     page, pageSize int,
// ) ([]entities.Reception, int64, error) {
//     op := "repo.Reception.GetPatientReceptionsWithPagination"

//     query := r.db.Model(&entities.Reception{}).
//         Preload("Specialization").
//         Where("patient_id = ?", patientID)

//     var totalRecords int64
//     if err := query.Count(&totalRecords).Error; err != nil {
//         return nil, 0, errors.NewDBError(op, err)
//     }

//     if page > 0 && pageSize > 0 {
//         offset := (page - 1) * pageSize
//         query = query.Offset(offset).Limit(pageSize)
//     }

//     var receptions []entities.Reception
//     result := query.Order("created_at DESC").Find(&receptions)
//     if result.Error != nil {
//         return nil, 0, errors.NewDBError(op, result.Error)
//     }

//     return receptions, totalRecords, nil
// }

// // GetReceptionWithFullData - получение приема с полными данными
// func (r *ReceptionRepository) GetReceptionWithFullData(receptionID uint) (*entities.Reception, error) {
//     op := "repo.Reception.GetReceptionWithFullData"

//     var reception entities.Reception
//     err := r.db.Preload("Specialization").
//         Preload("Patient").
//         Where("id = ?", receptionID).
//         First(&reception).Error

//     if err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             return nil, fmt.Errorf("%s: reception not found", op)
//         }
//         return nil, errors.NewDBError(op, err)
//     }

//     return &reception, nil
// }

// // SearchReceptions - поиск приемов по ключевым словам
// func (r *ReceptionRepository) SearchReceptions(
//     patientID uint,
//     searchTerm string,
// ) ([]entities.Reception, error) {
//     op := "repo.Reception.SearchReceptions"

//     var receptions []entities.Reception
//     err := r.db.Preload("Specialization").
//         Where("patient_id = ? AND (specialization_data::text ILIKE ?)",
//             patientID, "%"+searchTerm+"%").
//         Order("created_at DESC").
//         Find(&receptions).Error

//     if err != nil {
//         return nil, errors.NewDBError(op, err)
//     }

//     return receptions, nil
// }

// // GetReceptionsByDateRange - получение приемов за период
// func (r *ReceptionRepository) GetReceptionsByDateRange(
//     patientID uint,
//     startDate, endDate time.Time,
// ) ([]entities.Reception, error) {
//     op := "repo.Reception.GetReceptionsByDateRange"

//     var receptions []entities.Reception
//     err := r.db.Preload("Specialization").
//         Where("patient_id = ? AND created_at BETWEEN ? AND ?",
//             patientID, startDate, endDate).
//         Order("created_at DESC").
//         Find(&receptions).Error

//     if err != nil {
//         return nil, errors.NewDBError(op, err)
//     }

//     return receptions, nil
// }

// // GetReceptionStatistics - получение статистики по приемам
// func (r *ReceptionRepository) GetReceptionStatistics(patientID uint) (*ReceptionStatistics, error) {
//     op := "repo.Reception.GetReceptionStatistics"

//     var stats ReceptionStatistics

//     // Общее количество приемов
//     if err := r.db.Model(&entities.Reception{}).
//         Where("patient_id = ?", patientID).
//         Count(&stats.TotalReceptions).Error; err != nil {
//         return nil, errors.NewDBError(op, err)
//     }

//     // Количество завершенных приемов
//     if err := r.db.Model(&entities.Reception{}).
//         Where("patient_id = ? AND is_completed = ?", patientID, true).
//         Count(&stats.CompletedReceptions).Error; err != nil {
//         return nil, errors.NewDBError(op, err)
//     }

//     // Количество незавершенных приемов
//     if err := r.db.Model(&entities.Reception{}).
//         Where("patient_id = ? AND is_completed = ?", patientID, false).
//         Count(&stats.PendingReceptions).Error; err != nil {
//         return nil, errors.NewDBError(op, err)
//     }

//     return &stats, nil
// }

// // UpdateReceptionStatus - обновление статуса приема
// func (r *ReceptionRepository) UpdateReceptionStatus(receptionID uint, isCompleted bool) error {
//     op := "repo.Reception.UpdateReceptionStatus"

//     result := r.db.Model(&entities.Reception{}).
//         Where("id = ?", receptionID).
//         Updates(map[string]interface{}{
//             "is_completed": isCompleted,
//             "updated_at":    time.Now(),
//         })

//     if result.Error != nil {
//         return errors.NewDBError(op, result.Error)
//     }

//     if result.RowsAffected == 0 {
//         return fmt.Errorf("%s: reception not found", op)
//     }

//     return nil
// }

// // GetRecentReceptions - получение последних приемов
// func (r *ReceptionRepository) GetRecentReceptions(patientID uint, limit int) ([]entities.Reception, error) {
//     op := "repo.Reception.GetRecentReceptions"

//     var receptions []entities.Reception
//     err := r.db.Preload("Specialization").
//         Where("patient_id = ?", patientID).
//         Order("created_at DESC").
//         Limit(limit).
//         Find(&receptions).Error

//     if err != nil {
//         return nil, errors.NewDBError(op, err)
//     }

//     return receptions, nil
// }

// type ReceptionStatistics struct {
//     TotalReceptions      int64 `json:"total_receptions"`
//     CompletedReceptions  int64 `json:"completed_receptions"`
//     PendingReceptions    int64 `json:"pending_receptions"`
// }
