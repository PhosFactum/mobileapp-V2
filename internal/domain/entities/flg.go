package entities

import (
	"time"
)

// FLG - Флюорографическое исследование
type FLG struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	PatientID      uint `gorm:"not null;index" json:"patient_id"`
	OrganizationID uint `gorm:"not null" json:"organization_id"`

	ExaminationDate time.Time `gorm:"not null" json:"examination_date"`
	Number          string    `gorm:"type:varchar(100);not null" json:"number"`
	Result          string    `gorm:"type:text;not null" json:"result"`
	AttachedImage   string    `gorm:"type:varchar(500)" json:"attached_image,omitempty"`

	Patient      Patient      `gorm:"foreignKey:PatientID" json:"-"`
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"-"`
}
