package entities

import (
	"time"
)

type SpecialistResponse struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	PatientID      uint   `gorm:"not null;index" json:"patient_id"`
	SpecialistType string `gorm:"not null;index" json:"specialist_type"` // neurologist, psychiatrist, etc.

	// Ссылка на пациента
	Patient Patient `gorm:"foreignKey:PatientID" json:"-"`

	// Данные заключения в JSONB
	ResponseData SpecializationDataDocument `gorm:"type:jsonb" json:"response_data"`
}

// Типы специалистов
const (
	SpecialistNeurologist             = "neurologist"
	SpecialistPsychiatrist            = "psychiatrist"
	SpecialistTherapist               = "therapist"
	SpecialistPsychiatristNarcologist = "psychiatrist_narcologist"
	SpecialistEKG                     = "ekg"
	SpecialistMedicalCommission       = "medical_commission"
)

// Нужно будет подумать, вероятно придётся иначе определить связь
