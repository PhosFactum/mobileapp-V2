package usecases

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type ConsentUsecase struct {
	repo interfaces.ConsentSignatureRepository
}

func NewConsentUsecase(repo interfaces.ConsentSignatureRepository) *ConsentUsecase {
	return &ConsentUsecase{repo: repo}
}

// SaveConsent сохраняет подпись пациента и ставит флаг согласия
func (u *ConsentUsecase) SaveConsent(patientID uint, signature []byte) *errors.AppError {
	if err := u.repo.SaveSignature(patientID, signature); err != nil {
		return errors.NewAppError(
			errors.InternalServerErrorCode,
			"Не удалось сохранить подпись",
			err,
			false,
		)
	}
	return nil
}

// GetSignature возвращает подпись пациента (если нужна для фронта/админки)
func (u *ConsentUsecase) GetSignature(patientID uint) ([]byte, *errors.AppError) {
	signature, err := u.repo.GetSignature(patientID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.NewAppError(
				http.StatusNotFound,
				"Signature not found",
				err,
				true,
			)
		}

		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Не удалось получить подпись",
			err,
			false,
		)
	}
	return signature, nil
}

// GetConsentByPatientID возвращает полную сущность согласия (если нужна дата, ID и т.д.)
func (u *ConsentUsecase) GetConsentByPatientID(patientID uint) (*entities.ConsentSignature, *errors.AppError) {
	op := "usecase.Consent.GetConsentByPatientID"

	consent, err := u.repo.GetByPatientID(patientID)
	if err != nil {
		return nil, errors.NewInternalError(op, "не удалось получить данные согласия", err)
	}

	return consent, nil
}
