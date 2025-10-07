package reception

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type contextKey string

const txContextKey contextKey = "db_transaction"

type ReceptionRepositoryImpl struct {
	db *gorm.DB
}

func NewReceptionRepository(db *gorm.DB) interfaces.ReceptionRepository {
	return &ReceptionRepositoryImpl{db: db}
}
