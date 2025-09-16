package vaccine

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type VaccineRepository struct {
	db *gorm.DB
}

func NewVaccineRepository(db *gorm.DB) interfaces.VaccineRepository {
	repo := &VaccineRepository{db: db}
	return repo
}
