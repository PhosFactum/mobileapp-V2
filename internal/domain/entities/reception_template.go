package entities

import "encoding/json"

type ReceptionTemplate struct {
	ID            uint            `gorm:"primarykey"`
	Code          string          `gorm:"unique;not null"`
	Schema        json.RawMessage `gorm:"type:jsonb"`
	SchemaVersion string          `gorm:"not null"` //00002

	SpecializationID uint            `gorm:"not null;index" json:"specialization_id"`
	Specialization   *Specialization `gorm:"foreignKey:SpecializationID" json:"specialization"`

	HarmPoints []HarmPoint `gorm:"many2many:harm_point_reception_templates;"`
}
