package entities

import "time"

// Модели данных
type Organization struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title string `gorm:"not null" json:"title" example:"Med_Clinic"`

	ManagerID uint    `gorm:"not null;index" json:"-"`
	Manager   Manager `gorm:"foreignKey:ManagerID" json:"manager"`

	PatientGroups []PatientGroup `gorm:"foreignKey:OrganizationID" json:"patient_groups"`
	Doctor        []Doctor       `gorm:"many2many:doctor_organizations" json:"-"`
}

// Manager представляет информацию о мэнэджере организации
type Manager struct {
	ID uint `gorm:"primarykey" json:"id"`

	FullName string `gorm:"not null" json:"full_name" example:"Иванов Иван Иванович"`
	Phone    string `gorm:"unique;not null" json:"phone" example:"+79991234567"`
}
