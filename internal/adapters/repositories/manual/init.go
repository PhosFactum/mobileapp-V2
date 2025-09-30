package manual

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type ManualRepositoryImpl struct {
	db *gorm.DB
}

func NewManualRepository(db *gorm.DB) interfaces.ManualRepository {
	return &ManualRepositoryImpl{db: db}
}
