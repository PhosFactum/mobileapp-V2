package usecases

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type FLGUsecaseImpl struct {
	repo interfaces.FLGRepository
}

func NewFLGUsecase(repo interfaces.FLGRepository) *FLGUsecaseImpl {
	return &FLGUsecaseImpl{repo: repo}
}

func (u *FLGUsecaseImpl) CreateFLG(req *models.FLGCreateRequest) (models.FLGResponse, *errors.AppError) {
	flg := entities.FLG{
		PatientID:       req.PatientID,
		OrganizationID:  req.OrganizationID,
		ExaminationDate: req.ExaminationDate,
		Number:          req.Number,
		Result:          req.Result,
		AttachedImage:   req.AttachedImage,
	}

	id, err := u.repo.CreateFLG(flg)
	if err != nil {
		return models.FLGResponse{}, errors.NewInternalError("usecase.FLG.Create", "не удалось создать ФЛГ запись", err)
	}

	return models.FLGResponse{
		ID:              id,
		PatientID:       req.PatientID,
		OrganizationID:  req.OrganizationID,
		ExaminationDate: req.ExaminationDate,
		Number:          req.Number,
		Result:          req.Result,
		AttachedImage:   req.AttachedImage,
	}, nil
}

func (u *FLGUsecaseImpl) UpdateFLG(flgID uint, req *models.FLGUpdateRequest) (models.FLGResponse, *errors.AppError) {
	updateMap := make(map[string]interface{})

	if !req.ExaminationDate.IsZero() {
		updateMap["examination_date"] = req.ExaminationDate
	}
	if req.Number != "" {
		updateMap["number"] = req.Number
	}
	if req.Result != "" {
		updateMap["result"] = req.Result
	}
	if req.AttachedImage != "" {
		updateMap["attached_image"] = req.AttachedImage
	}
	if req.OrganizationID != 0 {
		updateMap["organization_id"] = req.OrganizationID
	}

	updated, err := u.repo.UpdateFLG(flgID, updateMap)
	if err != nil {
		return models.FLGResponse{}, errors.NewInternalError("usecase.FLG.Update", "не удалось обновить ФЛГ", err)
	}

	return models.FLGResponse{
		ID:              updated.ID,
		PatientID:       updated.PatientID,
		OrganizationID:  updated.OrganizationID,
		ExaminationDate: updated.ExaminationDate,
		Number:          updated.Number,
		Result:          updated.Result,
		AttachedImage:   updated.AttachedImage,
		CreatedAt:       updated.CreatedAt,
		UpdatedAt:       updated.UpdatedAt,
	}, nil
}
