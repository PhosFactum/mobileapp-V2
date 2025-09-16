package analysis

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type AnalysisRepository struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) interfaces.AnalysisRepository {
	repo := &AnalysisRepository{db: db}
	return repo
}
