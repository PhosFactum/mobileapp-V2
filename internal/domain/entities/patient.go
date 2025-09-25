package entities

import (
	"time"
)

// entities/patient.go

type Patient struct {
	ID              uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	FullName        string    `gorm:"not null" json:"full_name" example:"Иванов Иван Иванович"`
	BirthDate       time.Time `gorm:"not null" json:"birth_date" example:"1980-05-15T00:00:00Z"`
	IsMale          bool      `gorm:"not null" json:"is_male" example:"true"`
	Position        string    `gorm:"not null" json:"position" example:"Прогер"`
	Division        string    `gorm:"not null" json:"division" example:"Прогер"`
	ExaminationType string    `gorm:"not null" json:"examination_type_id"`
	ExaminationView string    `gorm:"not null" json:"examination_view_id"`

	PatientGroupID uint `gorm:"not null;index" json:"patient_group_id"`

	HarmPointID uint       `gorm:"not null;index" json:"harm_point_id"`
	HarmPoint   *HarmPoint `gorm:"foreignKey:HarmPointID" json:"harm_point,omitempty"`

	PersonalInfoID uint          `gorm:"not null;index" json:"personal_info_id"`
	PersonalInfo   *PersonalInfo `gorm:"foreignKey:PersonalInfoID" json:"personal_info,omitempty"`

	ContactInfoID uint         `gorm:"not null;index" json:"contact_info_id"`
	ContactInfo   *ContactInfo `gorm:"foreignKey:ContactInfoID" json:"contact_info,omitempty"`

	FlgID *uint `gorm:"index" json:"flg_id"`
	Flg   *Flg  `gorm:"foreignKey:FlgID" json:"flg,omitempty"`

	AnalysisOrderID uint           `gorm:"not null;index" json:"analysis_order_id"`
	AnalysisOrder   *AnalysisOrder `gorm:"foreignKey:AnalysisOrderID" json:"analysis_order,omitempty"`

	Statistics *PatientStatistics `gorm:"foreignKey:PatientID" json:"statistics,omitempty"`

	Vaccines           []Vaccine           `gorm:"foreignKey:PatientID" json:"vaccines,omitempty"`
	VaccineRefusals    []VaccineRefusal    `gorm:"foreignKey:PatientID" json:"vaccine_refusals,omitempty"`
	VaccineWithdrawals []VaccineWithdrawal `gorm:"foreignKey:PatientID" json:"vaccine_withdrawals,omitempty"`
	Titers             []Titr              `gorm:"foreignKey:PatientID" json:"titers,omitempty"`
	Receptions         []Reception         `gorm:"foreignKey:PatientID" json:"receptions,omitempty"`

	Specializations []Specialization `gorm:"many2many:patients_specializations;" json:"specializations,omitempty"`
}

type HarmPoint struct {
	ID    uint    `gorm:"primarykey" json:"id"`
	Value float32 `gorm:"not null;" json:"value"`

	Specializations []Specialization `gorm:"many2many:harm_points_specializations;" json:"-"`
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

	DocumentTypeID uint `gorm:"not null;index" json:"document_type_id" example:"1"`
}

type PatientStatistics struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"updated_at"`

	PatientID              uint  `gorm:"not null;uniqueIndex" json:"-"`
	TotalReceptions        int64 `gorm:"not null;default:0" json:"total_receptions"`
	CompletedReceptions    int64 `gorm:"not null;default:0" json:"completed_receptions"`
	TotalAnalysisOrders    int64 `gorm:"not null;default:0" json:"total_analysis_orders"`
	CompletedAnalysisItems int64 `gorm:"not null;default:0" json:"completed_analysis_items"`
}
