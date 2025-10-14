package patient_group

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type PatientGroupRepositoryImpl struct {
	db *gorm.DB
}

func NewPatientGroupRepository(db *gorm.DB) interfaces.PatientGroupRepository {
	return &PatientGroupRepositoryImpl{db: db}
}
