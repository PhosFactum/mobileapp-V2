package entities

import (
	"time"
)

// Doctor представляет информацию о враче
type Doctor struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FullName     string `gorm:"not null" json:"full_name" example:"Иванов Иван Иванович"`
	Phone        string `gorm:"unique;not null" json:"phone" example:"+79991234567"`
	PasswordHash string `gorm:"not null" json:"-"`

	Specializations []Specialization `gorm:"many2many:doctor_specializations" json:"-"`
	Organizations   []Organization   `gorm:"many2many:doctor_organizations" json:"-"`
}
