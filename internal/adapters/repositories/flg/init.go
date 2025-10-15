package flg

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/base"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type FlgRepositoryImpl struct {
	*base.BaseRepository
}

func NewFlgRepository(db *gorm.DB) interfaces.FlgRepository {
	return &FlgRepositoryImpl{
		BaseRepository: base.NewBaseRepository(db),
	}
}
