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

func (u *PatientGroupUsecase) GetPatientGroupsByOrganization(search string, page, perPage int,
) (*models.FilterResponse[[]models.PatientGroupShortResponse], error) {
	// Получаем данные из репозитория
	patientGroups, total, err := u.repo.GetByCodeOrOrgTitle(search, page, perPage)
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
