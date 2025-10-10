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

// GetTemplateSchemaByID возвращает только schema по template_id
func (r *ReceptionRepositoryImpl) GetTemplateSchemaByID(ctx context.Context, templateID uint) (json.RawMessage, error) {
	op := "repo.ReceptionTemplate.GetTemplateSchemaByID"

	var schema json.RawMessage
	err := r.GetDB(ctx).WithContext(ctx).
		Model(&entities.ReceptionTemplate{}).
		Select("schema").
		Where("id = ?", templateID).
		Scan(&schema).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("Reception Template Not Found")
		}
		return nil, errors.NewDBError(op, err)
	}
	return schema, nil
}

// CreateReception создаёт приём
func (r *ReceptionRepositoryImpl) CreateReception(ctx context.Context, reception *entities.Reception) error {
	op := "repo.Reception.CreateReception"
	if err := r.GetDB(ctx).WithContext(ctx).Create(reception).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}
