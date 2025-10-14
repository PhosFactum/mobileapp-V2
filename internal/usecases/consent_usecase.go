package usecases

import (
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
	op := "usecase.Consent.SaveConsent"

	if err := u.repo.SaveSignature(patientID, signature); err != nil {
		return errors.NewInternalError(op, "не удалось сохранить подпись", err)
	}

	return nil
}

// GetSignature возвращает подпись пациента (если нужна для фронта/админки)
func (u *ConsentUsecase) GetSignature(patientID uint) ([]byte, *errors.AppError) {
	op := "usecase.Consent.GetSignature"

	signature, err := u.repo.GetSignature(patientID)
	if err != nil {
		return nil, errors.NewInternalError(op, "не удалось получить подпись", err)
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
