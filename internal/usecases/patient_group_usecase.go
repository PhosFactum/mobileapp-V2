package usecases

import (
	"math"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type PatientGroupUsecase struct {
	repo interfaces.PatientGroupRepository
}

func NewPatientGroupUsecase(repo interfaces.PatientGroupRepository) interfaces.PatientGroupUseCase {
	return &PatientGroupUsecase{
		repo: repo}
}

func (u *PatientGroupUsecase) GetPatientGroupsByCodeOrOrgTitle(search string, page, perPage int,
) (*models.FilterResponse[[]models.PatientGroupShortResponse], error) {
	// Получаем данные из репозитория
	patientGroups, total, err := u.repo.GetPatientGroupsByCodeOrOrgTitle(search, page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get pateintGroups",
			err,
			false,
		)
	}

	// Преобразуем в DTO
	response := make([]models.PatientGroupShortResponse, len(patientGroups))
	for i, patient_group := range patientGroups {
		response[i] = models.PatientGroupShortResponse{
			ID:                patient_group.ID,
			Code:              patient_group.Code,
			OrganizationTitle: patient_group.Organization.Title,
			CreatedAt:         patient_group.CreatedAt.Format(time.RFC3339),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.FilterResponse[[]models.PatientGroupShortResponse]{
		Hits:        response,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalHits:   int(total),
		HitsPerPage: perPage,
	}, nil
}

func (u *PatientGroupUsecase) GetPatientGroupsByOrganizationID(orgID uint, page, perPage int) (*models.FilterResponse[[]models.PatientGroupWithPatientsResponse], error) {
	// Получаем данные из репозитория с пациентами
	patientGroups, total, err := u.repo.GetPatientGroupsWithPatientsByOrganizationID(orgID, page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get patient groups with patients",
			err,
			false,
		)
	}

	// Преобразуем в DTO с пациентами
	response := make([]models.PatientGroupWithPatientsResponse, len(patientGroups))
	for i, group := range patientGroups {
		// Преобразуем пациентов в ShortPatientResponse
		patients := make([]models.ShortPatientResponse, len(group.Patient))
		for j, patient := range group.Patient {
			patients[j] = models.ShortPatientResponse{
				ID:        patient.ID,
				FullName:  patient.FullName,
				BirthDate: patient.BirthDate,
				IsMale:    patient.IsMale,
			}
		}

		response[i] = models.PatientGroupWithPatientsResponse{
			ID:           group.ID,
			Code:         group.Code,
			Organization: group.Organization.Title,
			CreatedAt:    group.CreatedAt.Format(time.RFC3339),
			Patients:     patients,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.FilterResponse[[]models.PatientGroupWithPatientsResponse]{
		Hits:        response,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalHits:   int(total),
		HitsPerPage: perPage,
	}, nil
}
