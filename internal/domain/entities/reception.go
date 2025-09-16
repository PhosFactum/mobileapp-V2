package entities

import (
	"time"

	"github.com/jackc/pgtype"
)

// Reception заключения врачей
type Reception struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IsCompleted bool `gorm:"default:false" json:"is_completed"`

	PatientID uint    `gorm:"not null;index" json:"patient_id" example:"1"`
	Patient   Patient `gorm:"foreignKey:PatientID" json:"-"`

	// УНИКАЛЬНОЕ ОГРАНИЧЕНИЕ: у пациента может быть только одно заключение на специализацию
	// Добавляем композитный уникальный индекс
	// _ struct{} `gorm:"uniqueIndex:idx_patient_specialization"`

	// Связь со специализацией
	SpecializationID uint            `gorm:"not null;index" json:"specialization_id"`
	Specialization   *Specialization `gorm:"foreignKey:SpecializationID" json:"specialization"`

	SpecializationData pgtype.JSONB `gorm:"type:jsonb" json:"specialization_data"` // Данные как приходят из информационной системы
	CustomFieldsSchema []byte       `gorm:"type:jsonb" json:"-"`                   // Данные валидации
}

type SpecializationDataDocument struct {
	DocumentType string        `json:"document_type"`
	Fields       []CustomField `json:"fields"`
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

// Вспомогательная функция для создания указателя на int
func intPtr(i int) *int {
	return &i
}

// JSONB — фиктивный тип, чтобы Swagger знал, как дескриптить pgtype.JSONB
// @name JSONB
type JSONB map[string]interface{}
