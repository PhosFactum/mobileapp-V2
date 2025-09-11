package entities

import "time"

// Модели данных
type Organization struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title string `gorm:"not null" json:"title" example:"Med_Clinic"`

	ManagerID uint    `gorm:"not null;index" json:"-"`
	Manager   Manager `gorm:"foreignKey:ManagerID" json:"Manager"`

	PatientGroups []PatientGroup `gorm:"foreignKey:OrganizationID" json:"receptions"`
}
