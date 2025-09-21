package usecases

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

type PatientUsecase struct {
	repo          interfaces.PatientRepository
	contactRepo   interfaces.ContactInfoRepository
	personalRepo  interfaces.PersonalInfoRepository
	FilterBuilder interfaces.FilterBuilderService
}

func NewPatientUsecase(repo interfaces.PatientRepository, contactRepo interfaces.ContactInfoRepository, personalRepo interfaces.PersonalInfoRepository, s interfaces.Service) interfaces.PatientUsecase {
	return &PatientUsecase{
		repo:          repo,
		contactRepo:   contactRepo,
		personalRepo:  personalRepo,
		FilterBuilder: s}
}

// CreatePatient - создание пациента
func (u *PatientUsecase) CreatePatient(patientData *models.CreatePatientData) (*entities.Patient, error) {
	// // 1. Валидация входных данных
	// if err := u.validateCreatePatientData(patientData); err != nil {
	//     return nil, errors.NewValidationError(op, err)
	// }

	// // 2. Проверяем существование обязательных сущностей
	// if err := u.validateRequiredEntities(patientData); err != nil {
	//     return nil, errors.NewValidationError(op, err)
	// }

	// 3. Создаем пациента через репозиторий
	patient, err := u.repo.CreatePatient(patientData)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to create patient", err, true)
	}

	return patient, nil
}

func (u *PatientUsecase) GetPatientByID(id uint) (entities.Patient, *errors.AppError) {
	patient, err := u.repo.GetPatientByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, errors.ErrNotFound) {
			return entities.Patient{}, errors.NewAppError(
				http.StatusNotFound,
				"Пациент не найден",
				nil,
				true,
			)
		}

		return entities.Patient{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			errors.InternalServerError,
			err,
			false,
		)
	}
	return patient, nil
}

func (u *PatientUsecase) UpdatePatient(input *models.UpdatePatientRequest) (entities.Patient, *errors.AppError) {
	parsedTime, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		fmt.Println("Ошибка парсинга даты:", err)
		return entities.Patient{}, errors.NewAppError(errors.InvalidDataCode, "Ошибка парсинга даты:", err, false)
	}

	updateMap := map[string]interface{}{
		"id":         input.ID,
		"birth_date": parsedTime,
		"full_name":  input.FullName,
		"updated_at": time.Now(),
	}

	updatedPatientId, err := u.repo.UpdatePatient(input.ID, updateMap)
	if err != nil {
		return entities.Patient{}, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	updatedPatient, err := u.repo.GetPatientByID(updatedPatientId)
	if err != nil {
		return entities.Patient{}, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	return updatedPatient, nil

}

func (u *PatientUsecase) DeletePatient(id uint) *errors.AppError {
	if err := u.repo.DeletePatient(id); err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "удаление пациента", err, false)
	}
	return nil
}

// func (u *PatientUsecase) GetAllPatients(page, count int, filter string, order string) (models.FilterResponse[[]models.ShortPatientResponse], *errors.AppError) {
// 	var queryFilter string
// 	var queryOrder string
// 	var parameters []interface{}
// 	empty := models.FilterResponse[[]models.ShortPatientResponse]{}

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
// 	patients, totalRows, err := u.repo.GetAllPatients(page, count, queryFilter, queryOrder, parameters)
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

// 	var resp_models []models.ShortPatientResponse
// 	for _, patient := range patients {
// 		model := mapPatientEntityToModel(patient)
// 		resp_models = append(resp_models, model)
// 	}

// 	return models.FilterResponse[[]models.ShortPatientResponse]{
// 		Hits:        resp_models,
// 		CurrentPage: page,
// 		HitsPerPage: len(resp_models),
// 		TotalHits:   int(totalRows),
// 		TotalPages:  totalPages,
// 	}, nil
// }

func (u *PatientUsecase) GetPatientsByGroup(page, perPage int, group_id uint,
) (models.FilterResponse[[]entities.Patient], *errors.AppError) {

	empty := models.FilterResponse[[]entities.Patient]{}

	// Получение пациентов
	patients, totalRows, err := u.repo.GetPatientsByGroup(page, perPage, group_id)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to get patients", err, true)
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(perPage)))

	return models.FilterResponse[[]entities.Patient]{
		Hits:        patients,
		CurrentPage: page,
		HitsPerPage: len(patients),
		TotalHits:   int(totalRows),
		TotalPages:  totalPages,
	}, nil
}

func mapPatientEntityToModel(entity entities.Patient) models.ShortPatientResponse {
	return models.ShortPatientResponse{
		ID:        entity.ID,
		FullName:  entity.FullName,
		BirthDate: entity.BirthDate,
		IsMale:    entity.IsMale,
	}
}
