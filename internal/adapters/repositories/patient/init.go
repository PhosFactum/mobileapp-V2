package patient

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/base"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type contextKey string

const txContextKey contextKey = "db_transaction"

type PatientRepositoryImpl struct {
	*base.BaseRepository
}

func NewPatientRepository(db *gorm.DB) interfaces.PatientRepository {
	return &PatientRepositoryImpl{
		BaseRepository: base.NewBaseRepository(db),
	}
}
