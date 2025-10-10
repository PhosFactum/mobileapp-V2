package entities

import "time"

type PatientGroup struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Code string `gorm:"not null" json:"code" example:"94928490"`

	OrganizationID uint          `gorm:"not null;index" json:"-"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID" json:"Organization"`

	Patients []Patient `gorm:"foreignKey:PatientGroupID" json:"patients,omitempty"`
}
