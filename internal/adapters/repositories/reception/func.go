package reception

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

// getDB извлекает транзакцию из контекста или возвращает основное подключение
func (r *ReceptionRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return r.db
}

func (r *ReceptionRepositoryImpl) GetReceptionTemplatesByHarmPointID(ctx context.Context, harmPointID uint) ([]entities.ReceptionTemplate, error) {
	op := "repo.Patient.GetReceptionTemplatesByHarmPointID"

	var harmPoint entities.HarmPoint
	if err := r.getDB(ctx).WithContext(ctx).Preload("ReceptionTemplates").First(&harmPoint, harmPointID).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return harmPoint.ReceptionTemplates, nil
}

func (r *ReceptionRepositoryImpl) GetReceptionTemplatesByCodes(ctx context.Context, codes []string) ([]entities.ReceptionTemplate, error) {
	if len(codes) == 0 {
		return []entities.ReceptionTemplate{}, nil
	}
	var templates []entities.ReceptionTemplate
	if err := r.getDB(ctx).WithContext(ctx).Where("code IN ?", codes).Find(&templates).Error; err != nil {
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
	if err := r.getDB(ctx).WithContext(ctx).Create(&receptions).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}
