package usecases

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"

	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type ManualUseCase struct {
	repo interfaces.ManualRepository
}

func NewManualUseCase(repo interfaces.ManualRepository) interfaces.ManualUseCase {
	return &ManualUseCase{repo: repo}
}

func (u *ManualUseCase) GetAllManuals(ctx context.Context) ([]models.ManualResponse, *errors.AppError) {
	manuals, err := u.repo.GetAllManuals(ctx)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get manual entries",
			err,
			false,
		)
	}

	// Маппинг в ответную модель
	response := make([]models.ManualResponse, len(manuals))
	for i, m := range manuals {
		response[i] = models.ManualResponse{
			ID:    m.ID,
			Type:  string(m.Type),
			Value: m.Value,
		}
	}

	return response, nil
}
