package entities

import (
	"encoding/json"
	"time"
)

// Reception заключения врачей
type Reception struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IsCompleted bool `gorm:"default:false" json:"is_completed"`

	PatientID uint     `gorm:"not null;index" json:"patient_id" example:"1"`
	Patient   *Patient `gorm:"foreignKey:PatientID" json:"-"`

	// Связь со специализацией
	SpecializationID uint            `gorm:"not null;index" json:"specialization_id"`
	Specialization   *Specialization `gorm:"foreignKey:SpecializationID" json:"specialization"`

	// ✅ Объединённое поле: данные + схема валидации
	Data json.RawMessage `gorm:"type:jsonb" json:"data"`
}

type CustomField struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Required     bool        `json:"required"`
	Description  string      `json:"description"`
	Format       string      `json:"format,omitempty"`
	MinLength    *int        `json:"min_length,omitempty"`
	MaxLength    *int        `json:"max_length,omitempty"`
	MinValue     *int        `json:"min_value,omitempty"`
	MaxValue     *int        `json:"max_value,omitempty"`
	MinItems     *int        `json:"min_items,omitempty"`
	MaxItems     *int        `json:"max_items,omitempty"`
	Example      interface{} `json:"example,omitempty"`
	DefaultValue interface{} `json:"default_value,omitempty"`
	Value        interface{} `json:"value"`
	KeyFormat    string      `json:"key_format,omitempty"`
	ValueFormat  string      `json:"value_format,omitempty"`
}
