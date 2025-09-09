package entities

import "time"

// Модели данных
type Organization struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title string `gorm:"not null" json:"title" example:"Med_Clinic"`
	Code  string `gorm:"not null" json:"code" example:"89371943"`

	DoctorID uint    `gorm:"not null;index" json:"-"`
	Doctor   *Doctor `gorm:"foreignKey:DoctorID" json:"Doctor"`

	Patients []Patient `gorm:"foreignKey:OrganizationID" json:"patients"`
}
