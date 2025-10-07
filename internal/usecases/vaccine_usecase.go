package usecases

import (
	"context"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type VaccineUsecase struct {
	repo interfaces.VaccineRepository
}

func NewVaccineUsecase(repo interfaces.VaccineRepository) interfaces.VaccineUsecase {
	return &VaccineUsecase{
		repo: repo,
	}
}

func (u *VaccineUsecase) CreateVaccine(ctx context.Context, req *models.CreateVaccineRequest) (*entities.Vaccine, *errors.AppError) {
	op := "usecase.Vaccine.CreateVaccine"

	vaccine := &entities.Vaccine{
		Date:                req.Date,
		TitleID:             req.TitleID,
		PatientID:           req.PatientID,
		ResultID:            req.ResultID,
		MedicationID:        req.MedicationID,
		DoseID:              req.DoseID,
		NumberID:            req.NumberID,
		CertificateNumberID: req.CertificateNumberID,
		BodyPartID:          req.BodyPartID,
		MethodID:            req.MethodID,
		PlaceID:             req.PlaceID,
		CreatedAt:           time.Now(),
	}

	if err := u.repo.CreateVaccine(ctx, vaccine); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return vaccine, nil
}

func (u *VaccineUsecase) CreateVaccineRefusal(ctx context.Context, req *models.CreateVaccineRefusalRequest) (*entities.VaccineRefusal, *errors.AppError) {
	op := "usecase.Vaccine.CreateVaccineRefusal"

	refusal := &entities.VaccineRefusal{
		Date:      req.Date,
		TitleID:   req.TitleID,
		PatientID: req.PatientID,
		CreatedAt: time.Now(),
	}

	if err := u.repo.CreateVaccineRefusal(ctx, refusal); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return refusal, nil
}

func (u *VaccineUsecase) CreateVaccineWithdrawal(ctx context.Context, req *models.CreateVaccineWithdrawalRequest) (*entities.VaccineWithdrawal, *errors.AppError) {
	op := "usecase.Vaccine.CreateVaccineWithdrawal"

	withdrawal := &entities.VaccineWithdrawal{
		Date:      req.Date,
		TitleID:   req.TitleID,
		PatientID: req.PatientID,
		Num:       req.Num,
		CreatedAt: time.Now(),
	}

	if err := u.repo.CreateVaccineWithdrawal(ctx, withdrawal); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return withdrawal, nil
}

func (u *VaccineUsecase) CreateTitr(ctx context.Context, req *models.CreateTitrRequest) (*entities.Titr, *errors.AppError) {
	op := "usecase.Vaccine.CreateTitr"

	titr := &entities.Titr{
		Date:      req.Date,
		TitleID:   req.TitleID,
		PatientID: req.PatientID,
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}

	if err := u.repo.CreateTitr(ctx, titr); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return titr, nil
}
