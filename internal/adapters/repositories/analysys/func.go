package analysis

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *AnalysisRepositoryImpl) GetAnalysesByCodes(ctx context.Context, codes []string) ([]entities.Analysis, error) {
	if len(codes) == 0 {
		return []entities.Analysis{}, nil
	}
	var analyses []entities.Analysis
	if err := r.GetDB(ctx).WithContext(ctx).Where("code IN ?", codes).Find(&analyses).Error; err != nil {
		return nil, errors.NewDBError("repo.GetAnalysesByCodes", err)
	}
	return analyses, nil
}

func (a *AnalysisRepositoryImpl) GetAnalysesByHarmPointID(ctx context.Context, harmPointID uint) ([]entities.Analysis, error) {
	op := "repo.Analysis.GetAnalysesByHarmPointID"

	var harmPoint entities.HarmPoint
	if err := a.GetDB(ctx).WithContext(ctx).Preload("Analyses").First(&harmPoint, harmPointID).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return harmPoint.Analyses, nil
}

// repository/analysis_impl.go
func (r *AnalysisRepositoryImpl) GetAnalysisByID(ctx context.Context, id uint) (*entities.Analysis, error) {
	var a entities.Analysis
	err := r.GetDB(ctx).First(&a, id).Error
	return &a, err
}

func (r *AnalysisRepositoryImpl) GetAllAnalysisIDs(ctx context.Context) ([]uint, error) {
	var ids []uint
	err := r.GetDB(ctx).Model(&entities.Analysis{}).Pluck("id", &ids).Error
	return ids, err
}
