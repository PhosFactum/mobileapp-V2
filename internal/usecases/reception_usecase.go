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

// GetPatientReceptionStatuses возвращает список статусов заключений пациента
// func (u *ReceptionUsecase) GetPatientReceptionStatuses(patientID uint) ([]models.ReceptionStatus, *errors.AppError) {
// 	// Получаем приемы из репозитория
// 	receptions, err := u.repo.GetPatientReceptionsByPatientID(patientID)
// 	if err != nil {
// 		return nil, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"failed to fetch receptions for patient",
// 			err,
// 			true,
// 		)
// 	}

// 	// Преобразуем в DTO
// 	var result []models.ReceptionStatus
// 	for _, r := range receptions {
// 		statusText := "В процессе"
// 		if r.IsCompleted {
// 			statusText = "Завершено"
// 		}

// 		specTitle := "Неизвестно"
// 		if r.Specialization != nil {
// 			specTitle = r.Specialization.Title
// 		}

// 		result = append(result, models.ReceptionStatus{
// 			Specialization: specTitle,
// 			Status:         statusText,
// 		})
// 	}

// 	return result, nil
// }
