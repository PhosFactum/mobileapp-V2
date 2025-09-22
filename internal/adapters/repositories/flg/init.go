package flg

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type FLGRepositoryImpl struct {
	db *gorm.DB
}

func NewFLGRepository(db *gorm.DB) interfaces.FLGRepository {
	return &FLGRepositoryImpl{db: db}
}
