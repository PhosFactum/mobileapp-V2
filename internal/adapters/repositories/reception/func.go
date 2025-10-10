package reception

import (
	"context"
	"encoding/json"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

func (r *ReceptionRepositoryImpl) GetReceptionTemplatesByHarmPointID(ctx context.Context, harmPointID uint) ([]entities.ReceptionTemplate, error) {
	op := "repo.Patient.GetReceptionTemplatesByHarmPointID"

	var harmPoint entities.HarmPoint
	if err := r.GetDB(ctx).WithContext(ctx).Preload("ReceptionTemplates").First(&harmPoint, harmPointID).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return harmPoint.ReceptionTemplates, nil
}

func (r *ReceptionRepositoryImpl) GetReceptionTemplatesByCodes(ctx context.Context, codes []string) ([]entities.ReceptionTemplate, error) {
	if len(codes) == 0 {
		return []entities.ReceptionTemplate{}, nil
	}
	var templates []entities.ReceptionTemplate
	if err := r.GetDB(ctx).WithContext(ctx).Where("code IN ?", codes).Find(&templates).Error; err != nil {
		return nil, errors.NewDBError("repo.GetReceptionTemplatesByCodes", err)
	}
	return templates, nil
}

// CreateReceptions создаёт приёмы
func (r *ReceptionRepositoryImpl) CreateReceptions(ctx context.Context, receptions []entities.Reception) error {
	if len(receptions) == 0 {
		return nil
	}
	op := "repo.Patient.CreateReceptions"
	if err := r.GetDB(ctx).WithContext(ctx).Create(&receptions).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *ReceptionRepositoryImpl) GetTemplateByReceptionID(ctx context.Context, receptionID uint) (*entities.ReceptionTemplate, error) {
	op := "repo.Reception.GetTemplateByReceptionID"

	var template entities.ReceptionTemplate

	err := r.GetDB(ctx).
		Table("receptions r").
		Select("rt.*").
		Joins("JOIN reception_templates rt ON r.template_id = rt.id").
		Where("r.id = ?", receptionID).
		Scan(&template).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError(op)
		}
		return nil, errors.NewDBError(op, err)
	}

	return &template, nil
}

func (r *ReceptionRepositoryImpl) UpdateReceptionData(ctx context.Context, receptionID uint, data json.RawMessage) error {
	op := "repo.Reception.UpdateReceptionData"

	err := r.GetDB(ctx).
		Model(&entities.Reception{}).
		Where("id = ?", receptionID).
		Update("data", data).Error

	if err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}
