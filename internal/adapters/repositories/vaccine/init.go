package vaccine

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/base"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type VaccineRepositoryImpl struct {
	*base.BaseRepository
}

func NewVaccineRepository(db *gorm.DB) interfaces.VaccineRepository {
	return &VaccineRepositoryImpl{
		BaseRepository: base.NewBaseRepository(db),
	}
}
