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

	// ИСПРАВЛЕНО: Убираем обязательные поля, которые могут быть NULL
	ExaminationTypeID *uint            `gorm:"index" json:"-"` // Может быть NULL
	ExaminationType   *ExaminationType `gorm:"foreignKey:ExaminationTypeID" json:"-"`

	ExaminationViewID *uint            `gorm:"index" json:"-"` // Может быть NULL
	ExaminationView   *ExaminationView `gorm:"foreignKey:ExaminationViewID" json:"-"`

	HarmPointID *uint      `gorm:"index" json:"-"` // Может быть NULL
	HarmPoint   *HarmPoint `gorm:"foreignKey:HarmPointID" json:"-"`

	PersonalInfoID *uint         `gorm:"index" json:"-"` // Может быть NULL
	PersonalInfo   *PersonalInfo `gorm:"foreignKey:PersonalInfoID" json:"-"`

	ContactInfoID *uint        `gorm:"index" json:"-"` // Может быть NULL
	ContactInfo   *ContactInfo `gorm:"foreignKey:ContactInfoID" json:"-"`

	OrganizationID *uint         `gorm:"index" json:"organization_id,omitempty" example:"1"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID" json:"-"`

	FlgID *uint `gorm:"index" json:"-"`
	Flg   *Flg  `gorm:"foreignKey:FlgID" json:"flg"`

	Vaccines []Vaccine `gorm:"foreignKey:PatientID"`

	PatientGroups   []PatientGroup   `gorm:"many2many:patients_patient_groups;" json:"-"`
	Specializations []Specialization `gorm:"many2many:patients_specializations;" json:"-"`

	Statistics *PatientStatistics `gorm:"foreignKey:PatientID" json:"statistics,omitempty"`
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

	DocumentTypeID *uint         `gorm:"index" json:"document_type_iD" example:"1"`
	DocumentType   *DocumentType `gorm:"foreignKey:DocumentTypeID" json:"-"`
}

type DocumentType struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null;" json:"value"`
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

type PatientStatistics struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"updated_at"`

	PatientID uint `gorm:"not null;uniqueIndex" json:"patient_id"`

	// Статистика по приемам
	TotalReceptions     int64 `gorm:"not null;default:0" json:"total_receptions"`
	CompletedReceptions int64 `gorm:"not null;default:0" json:"completed_receptions"`

	// Статистика по анализам
	TotalAnalysisOrders    int64 `gorm:"not null;default:0" json:"total_analysis_orders"`
	CompletedAnalysisItems int64 `gorm:"not null;default:0" json:"completed_analysis_items"`
}
