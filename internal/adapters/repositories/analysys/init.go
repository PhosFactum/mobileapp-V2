package analysis

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type AnalysisRepositoryImpl struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) interfaces.AnalysisRepository {
	repo := &AnalysisRepositoryImpl{db: db}
	return repo
}
