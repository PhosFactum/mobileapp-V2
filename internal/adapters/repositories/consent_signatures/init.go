package consent_signatures

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type ConsentSignatureRepositoryImpl struct {
	db *gorm.DB
}

func NewConsentSignatureRepository(db *gorm.DB) interfaces.ConsentSignatureRepository {
	return &ConsentSignatureRepositoryImpl{db: db}
}
