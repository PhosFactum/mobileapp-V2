package reception

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/base"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type ReceptionRepositoryImpl struct {
	*base.BaseRepository
}

func NewReceptionRepository(db *gorm.DB) interfaces.ReceptionRepository {
	return &ReceptionRepositoryImpl{
		BaseRepository: base.NewBaseRepository(db),
	}
}
