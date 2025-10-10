package analysis

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/base"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type AnalysisOrderRepositoryImpl struct {
	*base.BaseRepository
}

func NewAnalysisOrderRepository(db *gorm.DB) interfaces.AnalysisOrderRepository {
	return &AnalysisOrderRepositoryImpl{
		BaseRepository: base.NewBaseRepository(db),
	}
}
