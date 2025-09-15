package entities

import (
	"time"

	"github.com/jackc/pgtype"
)

// ReceptionHospital представляет приёмы стационара и выезда
type Reception struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IsCompleted bool `gorm:"default:false" json:"is_completed"`

	PatientID uint    `gorm:"not null;index" json:"patient_id" example:"1"`
	Patient   Patient `gorm:"foreignKey:PatientID" json:"-"`

	// Связь со специализацией
	SpecializationID uint            `gorm:"not null;index" json:"specialization_id"`
	Specialization   *Specialization `gorm:"foreignKey:SpecializationID" json:"specialization"`

	// Специализированные данные
	SpecializationData pgtype.JSONB `gorm:"type:jsonb" json:"specialization_data" swaggertype:"object"`

	SpecializationDataDecoded interface{} `gorm:"-" json:"specialization_data_decoded"`
}
