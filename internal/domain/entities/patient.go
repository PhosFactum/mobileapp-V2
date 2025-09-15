package entities

import (
	"time"
)

// Patient представляет информацию о пациенте
type Patient struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FullName  string    `gorm:"not null" json:"full_name" example:"Иванов Иван Иванович"`
	BirthDate time.Time `gorm:"not null" json:"birth_date" example:"1980-05-15T00:00:00Z"`
	IsMale    bool      `gorm:"not null" json:"is_male" example:"true"`

	PersonalInfo   *PersonalInfo `gorm:"foreignKey:PersonalInfoID" json:"-"`
	PersonalInfoID *uint         `gorm:"default:null" json:"-"`

	ContactInfo   *ContactInfo `gorm:"foreignKey:ContactInfoID" json:"-"`
	ContactInfoID *uint        `gorm:"default:null" json:"-"`

	OrganizationID *uint         `gorm:"default:null;index" json:"organization_id,omitempty" example:"1"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID" json:"-"`

	AnalysisOrderID *uint          `gorm:"default:null;index" json:"analysis_order_id,omitempty" example:"1"`
	AnalysisOrder   *AnalysisOrder `gorm:"foreignKey:AnalysisOrderID" json:"-"`

	FlGID *uint `gorm:"not null;index" json:"-"`
	FLG   *FlG  `gorm:"foreignKey:FLGID" json:"flg"`

	Vaccines []Vaccine `gorm:"foreignKey:OrganizationID" json:"vaccines"`

	PatientGroup   []PatientGroup   `gorm:"many2many:patients_patient_groups; default:null;" json:"-"`
	Specialization []Specialization `gorm:"many2many:patients_specializations; default:null;" json:"-"`
}
