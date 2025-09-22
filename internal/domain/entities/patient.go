package entities

import (
	"time"
)

type Patient struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FullName  string    `gorm:"not null" json:"full_name" example:"Иванов Иван Иванович"`
	BirthDate time.Time `gorm:"not null" json:"birth_date" example:"1980-05-15T00:00:00Z"`
	IsMale    bool      `gorm:"not null" json:"is_male" example:"true"`
	Position  string    `gorm:"not null" json:"position" example:"Прогер"`
	Division  string    `gorm:"not null" json:"division" example:"Прогер"`

	ExaminationTypeID uint             `gorm:"not null;index" json:"-"`
	ExaminationType   *ExaminationType `gorm:"foreignKey:ExaminationTypeID" json:"-"`

	ExaminationViewID uint             `gorm:"not null;index" json:"-"`
	ExaminationView   *ExaminationView `gorm:"foreignKey:ExaminationViewID" json:"-"`

	HarmPointID uint       `gorm:"not null;index" json:"-"`
	HarmPoint   *HarmPoint `gorm:"foreignKey:HarmPointID" json:"-"`

	PersonalInfoID uint          `gorm:"not null;index" json:"-"`
	PersonalInfo   *PersonalInfo `gorm:"foreignKey:PersonalInfoID" json:"-"`

	ContactInfoID uint         `gorm:"not null;index" json:"-"`
	ContactInfo   *ContactInfo `gorm:"foreignKey:ContactInfoID" json:"-"`

	FlgID *uint `gorm:"column:flg_id;index" json:"-"`
	Flg   *Flg  `gorm:"foreignKey:FlgID" json:"flg"`

	OrganizationID uint          `gorm:"not null;index" json:"organization_id,omitempty" example:"1"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID" json:"-"`

	Vaccines []Vaccine `gorm:"foreignKey:PatientID"`

	PatientGroups   []PatientGroup   `gorm:"many2many:patients_patient_groups;" json:"-"`
	Specializations []Specialization `gorm:"many2many:patients_specializations;" json:"-"`

	Statistics *PatientStatistics `gorm:"foreignKey:PatientID" json:"statistics,omitempty"`
}

type PatientStatistics struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"updated_at"`

	PatientID uint `gorm:"not null;uniqueIndex" json:"patient_id"`

	TotalReceptions     int64 `gorm:"not null;default:0" json:"total_receptions"`
	CompletedReceptions int64 `gorm:"not null;default:0" json:"completed_receptions"`

	TotalAnalysisOrders    int64 `gorm:"not null;default:0" json:"total_analysis_orders"`
	CompletedAnalysisItems int64 `gorm:"not null;default:0" json:"completed_analysis_items"`
}

type ExaminationType struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null;" json:"value"`
}

type ExaminationView struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null;" json:"value"`
}

type HarmPoint struct {
	ID    uint    `gorm:"primarykey" json:"id"`
	Value float32 `gorm:"not null;" json:"value"`
}

type Flg struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	PatientID       uint      `gorm:"not null;index" json:"patient_id"`
	OrganizationID  uint      `gorm:"not null;index" json:"organization_id"`
	DoctorID        uint      `gorm:"not null;index" json:"doctor_id"`
	ExaminationDate time.Time `gorm:"not null" json:"examination_date"`
	Number          string    `gorm:"not null" json:"number" example:"ФЛГ-2025-001"`
	Result          string    `gorm:"not null" json:"result" example:"Годен"`
	AttachedImage   string    `json:"attached_image,omitempty"` // URL или base64
}

// ContactInfo представляет контактную информацию пациента
type ContactInfo struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// PersonalInfo представляет персональную информацию
type PersonalInfo struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DocNumber string `json:"doc_number" example:"4510 123456" rus:"Номер документа"`
	DocSeries string `json:"doc_series" example:"4510 123456" rus:"Серия документа"`
	SNILS     string `json:"snils" example:"123-456-789 00" rus:"СНИЛС"`
	OMS       string `json:"oms" example:"1234567890123456" rus:"Полис ОМС"`

	DocumentTypeID uint          `gorm:"not null;index" json:"document_type_iD" example:"1"`
	DocumentType   *DocumentType `gorm:"foreignKey:DocumentTypeID" json:"-"`
}

type DocumentType struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null;" json:"value"`
}
