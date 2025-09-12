package organization

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type OrganizationRepositoryImpl struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) interfaces.OrganizationRepository {
	return &OrganizationRepositoryImpl{db: db}
}
