// internal/domain/entities/consent_signature.go
package entities

import (
	"time"
)

type ConsentSignature struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PatientID uint      `gorm:"not null;index" json:"patient_id"`
	Signature []byte    `gorm:"type:bytea;not null" json:"-"` // Бинарные данные подписи
	Patient   Patient   `gorm:"foreignKey:PatientID" json:"-"`
}
