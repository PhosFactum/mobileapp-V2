package usecases

import (
	"math"

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

func (u *OrganizationUsecase) GetAllOrganizations(
	call_id uint,
	page, perPage int,
) (*models.FilterResponse[[]models.OrganizationShortResponse], error) {
	// Получаем данные из репозитория
	organizations, total, err := u.repo.GetAllOrganizations(page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get receptions",
			err,
			false,
		)
	}

	// Преобразуем в DTO
	response := make([]models.OrganizationShortResponse, len(organizations))
	for i, org := range organizations {
		response[i] = models.OrganizationShortResponse{
			ID:             org.ID,
			Title:          org.Title,
			Code:           org.Code,
			DoctorFullName: org.Doctor.FullName,
		}
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

func (u *OrganizationUsecase) GetByTitleOrCodeOrganizations(search string,
	call_id uint,
	page, perPage int,
) (*models.FilterResponse[[]models.OrganizationShortResponse], error) {
	// Получаем данные из репозитория
	organizations, total, err := u.repo.GetByTitleOrCodeOrganizations(search, page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get receptions",
			err,
			false,
		)
	}

	// Преобразуем в DTO
	response := make([]models.OrganizationShortResponse, len(organizations))
	for i, org := range organizations {
		response[i] = models.OrganizationShortResponse{
			ID:             org.ID,
			Title:          org.Title,
			Code:           org.Code,
			DoctorFullName: org.Doctor.FullName,
		}
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
