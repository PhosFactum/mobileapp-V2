package usecases

import (
	"math"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type OrganizationUsecase struct {
	repo interfaces.OrganizationRepository
}

func NewOrganizationUsecase(repo interfaces.OrganizationRepository) interfaces.OrganizationUseCase {
	return &OrganizationUsecase{
		repo: repo}
}

func (u *OrganizationUsecase) GetAllDoctorOrganizations(
	doctorID uint,
	search string,
	page, perPage int,
) (*models.FilterResponse[[]models.OrganizationShortResponse], *errors.AppError) {

	organizations, total, err := u.repo.GetAllDoctorOrganizations(doctorID, search, page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get organizations",
			err,
			false,
		)
	}

	response := make([]models.OrganizationShortResponse, len(organizations))
	for i, org := range organizations {
		response[i] = u.mapOrganizationShort(org)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.FilterResponse[[]models.OrganizationShortResponse]{
		Hits:        response,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalHits:   int(total),
		HitsPerPage: perPage,
	}, nil
}

// ✅ Внутренний маппер: Organization → OrganizationShortResponse
func (u *OrganizationUsecase) mapOrganizationShort(org entities.Organization) models.OrganizationShortResponse {
	return models.OrganizationShortResponse{
		ID:      org.ID,
		Title:   org.Title,
		Manager: u.mapManager(org.Manager),
	}
}

// ✅ Внутренний маппер: Manager → ManagerResponse
func (u *OrganizationUsecase) mapManager(manager entities.Manager) models.ManagerResponse {
	return models.ManagerResponse{
		FullName: manager.FullName,
		Phone:    manager.Phone,
	}
}
