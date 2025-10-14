package analysis

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/base"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type AnalysisRepositoryImpl struct {
	*base.BaseRepository
}

func NewAnalysisRepository(db *gorm.DB) interfaces.AnalysisRepository {
	return &AnalysisRepositoryImpl{
		BaseRepository: base.NewBaseRepository(db),
	}
}
