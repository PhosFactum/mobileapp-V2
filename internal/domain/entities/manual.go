package entities

import "time"

// Общая сущность для справочников ключ-значение
type ReferenceType string

const (
	RefTypeVaccineTitle             ReferenceType = "vaccine_title"
	RefTypeVaccineMedication        ReferenceType = "vaccine_medication"
	RefTypeVaccineDose              ReferenceType = "vaccine_dose"
	RefTypeVaccineNumber            ReferenceType = "vaccine_number"
	RefTypeVaccineCertificateNumber ReferenceType = "vaccine_certificate_number"
	RefTypeVaccineBodyPart          ReferenceType = "vaccine_body_part"
	RefTypeVaccineMethod            ReferenceType = "vaccine_method"
	RefTypeVaccinePlace             ReferenceType = "vaccine_place"
	RefTypePatientExaminationType   ReferenceType = "patient_examination_type"
	RefTypePatientExaminationView   ReferenceType = "patient_examination_view"
	RefTypePersonalDocumentType     ReferenceType = "personal_document_type"
	RefTypeMandatoryReception       ReferenceType = "mandatory_reception"
	RefTypeMandatoryAnalysis        ReferenceType = "mandatory_analysis"
)

type Manual struct {
	ID        uint          `gorm:"primarykey" json:"id"`
	Type      ReferenceType `gorm:"not null;index" json:"type"`
	Value     string        `gorm:"not null" json:"value"` // всегда строка, но парсится по ValueType
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
