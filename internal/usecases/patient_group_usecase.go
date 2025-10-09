package usecases

import (
	"math"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
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

func (u *PatientGroupUsecase) GetPatientGroupsByDoctorID(
	doctorID uint,
	search string,
	page, perPage int,
) (*models.FilterResponse[[]models.PatientGroupShortResponse], *errors.AppError) {

	patientGroups, total, err := u.repo.GetPatientGroupsByDoctorID(doctorID, search, page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get patient groups for doctor",
			err,
			false,
		)
	}

	response := make([]models.PatientGroupShortResponse, len(patientGroups))
	for i, pg := range patientGroups {
		response[i] = u.mapPatientGroupShort(pg)
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

// GetPatientGroupsByOrganizationID возвращает группы пациентов организации с поиском по коду
func (u *PatientGroupUsecase) GetPatientGroupsByOrganizationID(orgID uint, search string, page, perPage int) (*models.FilterResponse[[]models.PatientGroupShortResponse], *errors.AppError) {

	patientGroups, total, err := u.repo.GetPatientGroupsByOrganizationID(orgID, search, page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get patient groups for organization",
			err,
			false,
		)
	}

	response := make([]models.PatientGroupShortResponse, len(patientGroups))
	for i, pg := range patientGroups {
		response[i] = u.mapPatientGroupShort(pg)
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

// ✅ Внутренний маппер: PatientGroup → PatientGroupShortResponse
func (u *PatientGroupUsecase) mapPatientGroupShort(pg entities.PatientGroup) models.PatientGroupShortResponse {
	return models.PatientGroupShortResponse{
		ID:                pg.ID,
		Code:              pg.Code,
		OrganizationTitle: pg.Organization.Title, // ← Убедись, что Organization предзагружена!
		CreatedAt:         pg.CreatedAt.Format(time.RFC3339),
	}
}
