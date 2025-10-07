package analysis

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

// getDB извлекает транзакцию из контекста или возвращает основное подключение
func (r *AnalysisRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return r.db
}

func (r *AnalysisRepositoryImpl) CreateAnalysisItems(ctx context.Context, items []entities.AnalysisOrderItem) error {
	op := "repo.Patient.CreateAnalysisItems"
	if err := r.getDB(ctx).WithContext(ctx).Create(&items).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreateAnalysisOrder создаёт направление на анализы
func (r *AnalysisRepositoryImpl) CreateAnalysisOrder(ctx context.Context, order *entities.AnalysisOrder) error {
	op := "repo.Patient.CreateAnalysisOrder"
	if err := r.getDB(ctx).WithContext(ctx).Create(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// UpdateAnalysisOrder обновляет направление
func (r *AnalysisRepositoryImpl) UpdateAnalysisOrder(ctx context.Context, order *entities.AnalysisOrder) error {
	op := "repo.Patient.UpdateAnalysisOrder"
	if err := r.getDB(ctx).WithContext(ctx).Save(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *AnalysisRepositoryImpl) GetAnalysesByCodes(ctx context.Context, codes []string) ([]entities.Analysis, error) {
	if len(codes) == 0 {
		return []entities.Analysis{}, nil
	}
	var analyses []entities.Analysis
	if err := r.getDB(ctx).WithContext(ctx).Where("code IN ?", codes).Find(&analyses).Error; err != nil {
		return nil, errors.NewDBError("repo.GetAnalysesByCodes", err)
	}
	return analyses, nil
}

func (a *AnalysisRepositoryImpl) GetAnalysesByHarmPointID(ctx context.Context, harmPointID uint) ([]entities.Analysis, error) {
	op := "repo.Analysis.GetAnalysesByHarmPointID"

	var harmPoint entities.HarmPoint
	if err := a.getDB(ctx).WithContext(ctx).Preload("Analyses").First(&harmPoint, harmPointID).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return harmPoint.Analyses, nil
}
