package manual

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/base"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type ManualRepositoryImpl struct {
	*base.BaseRepository // ← ВСТРАИВАЕМ
}

func NewManualRepository(db *gorm.DB) interfaces.ManualRepository {
	return &ManualRepositoryImpl{
		BaseRepository: base.NewBaseRepository(db),
	}
}
