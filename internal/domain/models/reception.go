package models

import "encoding/json"

type ReceptionResponse struct {
	ID               uint                      `json:"id"`
	IsCompleted      bool                      `json:"is_completed"`
	SpecializationID uint                      `json:"specialization_id"`
	Specialization   SpecializationResponse    `json:"specialization,omitempty"`
	Template         ReceptionTemplateResponse `json:"template"`
	Data             json.RawMessage           `json:"data"`
}

type ReceptionTemplateResponse struct {
	ID     uint              `json:"id"`
	Code   string            `json:"code"`
	Fields []FieldDescriptor `json:"fields"`
}

type FieldDescriptor struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Type     string `json:"type"` // "string", "number", "integer", "boolean", "array"
	Required bool   `json:"required"`
	Tag      string `json:"tag"` // "input", "select", "checkbox", "textarea"

	// Для string
	MinLength *int    `json:"min_length,omitempty"`
	MaxLength *int    `json:"max_length,omitempty"`
	Pattern   *string `json:"pattern,omitempty"` // регулярное выражение

	// Для number/integer
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusive_minimum,omitempty"` // >, не >=
	ExclusiveMaximum *float64 `json:"exclusive_maximum,omitempty"` // <, не <=
	MultipleOf       *float64 `json:"multiple_of,omitempty"`       // кратность

	// Для array
	MinItems *int `json:"min_items,omitempty"`
	MaxItems *int `json:"max_items,omitempty"`
	// Items *FieldDescriptor // ← если нужны вложенные объекты (рекурсия)

	// Справочники
	Enum []string `json:"enum,omitempty"`

	// Форматы (для UI и валидации)
	Format *string `json:"format,omitempty"` // "email", "date", "phone"

	// Описание (для тултипов)
	Description *string `json:"description,omitempty"`
}
