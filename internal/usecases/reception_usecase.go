package usecases

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
)

type ReceptionUsecase struct {
	repo          interfaces.ReceptionRepository
	FilterBuilder interfaces.FilterBuilderService
}

func NewReceptionUsecase(repo interfaces.ReceptionRepository, s interfaces.Service) interfaces.ReceptionUsecase {
	return &ReceptionUsecase{
		repo:          repo,
		FilterBuilder: s}
}

// func (u *ReceptionUsecase) UpdateReceptionHospital(id uint, input *models.UpdateReceptionHospitalRequest) (models.ReceptionHospitalResponse, *errors.AppError) {

// 	recepHospUpdate := map[string]interface{}{
// 		"diagnosis":       input.Diagnosis,
// 		"recommendations": input.Recommendations,
// 	}

// 	if input.SpecializationData != nil {
// 		var specializationData pgtype.JSONB
// 		jsonData, err := json.Marshal(input.SpecializationData)
// 		if err != nil {
// 			return models.ReceptionHospitalResponse{}, errors.NewAppError(
// 				errors.InvalidDataCode,
// 				"failed to serialize specialization_data",
// 				err,
// 				true,
// 			)
// 		}

// 		if err := specializationData.Set(json.RawMessage(jsonData)); err != nil {
// 			return models.ReceptionHospitalResponse{}, errors.NewAppError(
// 				errors.InvalidDataCode,
// 				"failed to convert specialization_data to JSONB",
// 				err,
// 				true,
// 			)
// 		}

// 		recepHospUpdate["specialization_data"] = specializationData
// 	}

// 	if _, err := u.repo.UpdateReception(id, recepHospUpdate); err != nil {
// 		return models.ReceptionHospitalResponse{}, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"failed to update reception hospital data",
// 			err,
// 			true,
// 		)
// 	}

// 	reception, err := u.repo.GetReceptionByID(id)
// 	if err != nil {
// 		return models.ReceptionHospitalResponse{}, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"failed to get updated reception hospital data",
// 			err,
// 			true,
// 		)
// 	}
// 	return models.ReceptionHospitalResponse{
// 		Doctor: models.DoctorInfoResponse{
// 			FullName:       reception.Doctor.FullName,
// 			Specialization: reception.CachedSpecialization,
// 		},
// 		Patient: models.ShortPatientResponse{
// 			ID:         reception.Patient.ID,
// 			LastName:   reception.Patient.LastName,
// 			FirstName:  reception.Patient.FirstName,
// 			MiddleName: reception.Patient.MiddleName,
// 			BirthDate:  reception.Patient.BirthDate,
// 			IsMale:     reception.Patient.IsMale,
// 		},
// 		Diagnosis:       reception.Diagnosis,
// 		Recommendations: reception.Recommendations,
// 		Date:            reception.Date,
// 	}, nil
// }

// func (u *ReceptionUsecase) GetHospitalReceptionsByDoctorID(doc_id uint, page, count int, filter, order string) (models.FilterResponse[[]models.ReceptionHospitalResponse], *errors.AppError) {
// 	var queryFilter string
// 	var queryOrder string
// 	var parameters []interface{}

// 	empty := models.FilterResponse[[]models.ReceptionHospitalResponse]{}

// 	// Статические поля модели (имя таблицы/колонки и их типы)
// 	entityFields, err := getFieldTypes(entities.ReceptionHospital{})
// 	if err != nil {
// 		return empty, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, true)
// 	}

// 	// Парсим фильтр, если он передан
// 	if len(filter) > 0 {
// 		subQuery, params, err := u.FilterBuilder.ParseFilterString(filter, entityFields)
// 		if err != nil {
// 			return empty, errors.NewAppError(
// 				errors.InvalidDataCode,
// 				fmt.Sprintf("invalid filter syntax: %s", err.Error()),
// 				nil,
// 				false,
// 			)
// 		}
// 		queryFilter = subQuery
// 		parameters = params
// 	}

// 	// Сортировка по умолчанию по статусу и дате-времени
// 	if len(order) > 0 {
// 		subQuery, err := u.FilterBuilder.ParseOrderString(order, entityFields)
// 		if err != nil {
// 			return empty, errors.NewAppError(errors.InternalServerErrorCode, fmt.Sprintf("invalid order syntax: %s", err.Error()), nil, false)
// 		}
// 		queryOrder = subQuery + `,
// 		CASE status
// 			WHEN 'scheduled' THEN 0
// 			ELSE 1
// 		END, date ASC`
// 	} else {
// 		queryOrder = `CASE status
// 		WHEN 'scheduled' THEN 0
// 		ELSE 1
// 	END, date ASC`
// 	}

// 	// Получение пациентов
// 	receptions, totalRows, err := u.repo.GetAllHospitalReceptionsByDoctorID(doc_id, page, count, queryFilter, queryOrder, parameters)
// 	if err != nil {
// 		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to get patients", err, true)
// 	}

// 	var totalPages int
// 	if count == 0 {
// 		// Если count == 0, то пагинация отключена, и все записи возвращаются на одной странице
// 		totalPages = 1
// 		page = 1

// 	} else {
// 		// Вычисляем количество страниц с округлением вверх
// 		totalPages = int(math.Ceil(float64(totalRows) / float64(count)))
// 	}

// 	// Преобразуем в DTO
// 	response := make([]models.ReceptionHospitalResponse, len(receptions))
// 	for i, rec := range receptions {

// 		doctor := models.DoctorInfoResponse{
// 			FullName:       rec.Doctor.FullName,
// 			Specialization: rec.CachedSpecialization,
// 		}

// 		patient := models.ShortPatientResponse{
// 			ID:         rec.PatientID,
// 			LastName:   rec.Patient.LastName,
// 			FirstName:  rec.Patient.FirstName,
// 			MiddleName: rec.Patient.MiddleName,
// 			BirthDate:  rec.Patient.BirthDate,
// 			IsMale:     rec.Patient.IsMale,
// 		}

// 		response[i] = models.ReceptionHospitalResponse{
// 			ID:              rec.ID,
// 			Doctor:          doctor,
// 			Patient:         patient,
// 			Diagnosis:       rec.Diagnosis,
// 			Recommendations: rec.Recommendations,
// 			Status:          string(rec.Status),
// 			Date:            rec.Date,
// 		}
// 	}

// 	return models.FilterResponse[[]models.ReceptionHospitalResponse]{
// 		Hits:        response,
// 		CurrentPage: page,
// 		HitsPerPage: len(response),
// 		TotalHits:   int(totalRows),
// 		TotalPages:  totalPages,
// 	}, nil
// }

// func getStatusText(status entities.HospitalReceptionStatus) string {
// 	switch status {
// 	case entities.HospitalReceptionStatusScheduled:
// 		return "Запланирован"
// 	case entities.HospitalReceptionStatusCompleted:
// 		return "Завершен"
// 	case entities.HospitalReceptionStatusCancelled:
// 		return "Отменен"
// 	case entities.HospitalReceptionStatusNoShow:
// 		return "Не явился"
// 	default:
// 		return string(status)
// 	}
// }

// func (u *ReceptionHospitalUsecase) GetHospitalPatientsByDoctorID(
// 	doc_id uint,
// 	page, count int,
// 	filter, order string) (
// 	models.FilterResponse[[]entities.Patient], *errors.AppError) {

// 	var queryFilter string
// 	var queryOrder string
// 	var parameters []interface{}
// 	empty := models.FilterResponse[[]entities.Patient]{}

// 	if doc_id <= 0 {
// 		return empty, errors.NewAppError(errors.InternalServerErrorCode, "invalid doc_id", nil, true)
// 	}
// 	// Статические поля модели (имя таблицы/колонки и их типы)
// 	entityFields, err := getFieldTypes(entities.Patient{})
// 	if err != nil {
// 		return empty, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
// 	}

// 	// Парсим фильтр, если он передан
// 	if len(filter) > 0 {
// 		subQuery, params, err := u.FilterBuilder.ParseFilterString(filter, entityFields)
// 		if err != nil {
// 			return empty, errors.NewAppError(
// 				errors.InvalidDataCode,
// 				fmt.Sprintf("invalid filter syntax: %s", err.Error()),
// 				nil,
// 				false,
// 			)
// 		}
// 		queryFilter = subQuery
// 		parameters = params
// 	}

// 	if len(order) > 0 {
// 		subQuery, err := u.FilterBuilder.ParseOrderString(order, entityFields)
// 		if err != nil {
// 			return empty, errors.NewAppError(errors.InternalServerErrorCode, fmt.Sprintf("invalid order syntax: %s", err.Error()), nil, false)
// 		}
// 		queryOrder = subQuery
// 	}

// 	// Получение пациентов
// 	patients, totalRows, err := u.repo.GetAllPatientsFromHospitalByDoctorID(doc_id, page, count, queryFilter, queryOrder, parameters)
// 	if err != nil {
// 		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to get patients", err, true)
// 	}

// 	var totalPages int
// 	if count == 0 {
// 		// Если count == 0, то пагинация отключена, и все записи возвращаются на одной странице
// 		totalPages = 1
// 		page = 1

// 	} else {
// 		// Вычисляем количество страниц с округлением вверх
// 		totalPages = int(math.Ceil(float64(totalRows) / float64(count)))
// 	}

// 	return models.FilterResponse[[]entities.Patient]{
// 		Hits:        patients,
// 		CurrentPage: page,
// 		HitsPerPage: len(patients),
// 		TotalHits:   int(totalRows),
// 		TotalPages:  totalPages,
// 	}, nil
// }

// func (u *ReceptionHospitalUsecase) GetReceptionHospitalByID(
// 	hospID uint,
// ) (models.ReceptionFullResponse, error) {
// 	// Получаем данные из репозитория
// 	reception, err := u.repo.GetReceptionHospitalByID(hospID)
// 	if err != nil {
// 		return models.ReceptionFullResponse{}, fmt.Errorf("failed to get reception: %w", err)
// 	}

// 	doctor := models.DoctorShortResponse{
// 		ID:             reception.Doctor.ID,
// 		FullName:       reception.Doctor.FullName,
// 		Specialization: reception.CachedSpecialization,
// 	}

// 	// Формируем специализированные данные
// 	var specData interface{}
// 	if reception.SpecializationDataDecoded != nil {
// 		specData = reception.SpecializationDataDecoded
// 	} else if reception.SpecializationData.Status == pgtype.Present {
// 		// Если данные не декодированы, но есть в JSONB
// 		var rawData map[string]interface{}
// 		if err := reception.SpecializationData.AssignTo(&rawData); err == nil {
// 			specData = rawData
// 		}
// 	}

// 	response := models.ReceptionFullResponse{
// 		ID:                 reception.ID,
// 		Date:               reception.Date.Format("02.01.2006 15:04"),
// 		Status:             getStatusText(reception.Status),
// 		LastName:           reception.Patient.LastName,
// 		FirstName:          reception.Patient.FirstName,
// 		MiddleName:         reception.Patient.MiddleName,
// 		PatientID:          reception.Patient.ID,
// 		Diagnosis:          reception.Diagnosis,
// 		Address:            reception.Address,
// 		Doctor:             doctor,
// 		Recommendations:    reception.Recommendations,
// 		SpecializationData: specData,
// 	}

// 	return response, nil
// }

// func (u *ReceptionHospitalUsecase) UpdateReceptionHospitalStatus(id uint, newStatus string) (entities.ReceptionHospital, error) {
// 	empty := entities.ReceptionHospital{}

// 	// Обновление статуса в базе
// 	updateFields := map[string]interface{}{
// 		"status": newStatus,
// 	}

// 	if _, err := u.repo.UpdateReceptionHospital(id, updateFields); err != nil {
// 		return empty, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"failed to update reception hospital status",
// 			err,
// 			true,
// 		)
// 	}

// 	// Получение обновлённой записи
// 	reception, err := u.repo.GetReceptionHospitalByID(id)
// 	if err != nil {
// 		return empty, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"failed to get updated reception hospital record",
// 			err,
// 			true,
// 		)
// 	}

// 	return reception, nil

// }
