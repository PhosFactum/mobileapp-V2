package entities

import "time"

type Vaccine struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`

	Date        time.Time `json:"date"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`

	// Булевые флаги определяют подтип формы
	IsRefusal   bool `gorm:"default:false" json:"is_refusal"`
	IsExemption bool `gorm:"default:false" json:"is_exemption"`

	// Поле для титра
	TiterAmount *int `json:"titer_amount"`

	// Поле для медотвода
	MedWithdrawlNum *int `json:"med_withdrawl_num"`

	// Полная прививка
	Result *string `json:"result"`

	TitleID *uint  `json:"title_id,omitempty"` // Так же должно быть и для титра
	Title   *Title `gorm:"foreignKey:TitleID" json:"-"`

	MedicationID *uint       `json:"medication_id,omitempty"`
	Medication   *Medication `gorm:"foreignKey:MedicationID" json:"-"`

	DoseID *uint `json:"dose_id,omitempty"`
	Dose   *Dose `gorm:"foreignKey:DoseID" json:"-"`

	NumberID *uint   `json:"number_id,omitempty"`
	Number   *Number `gorm:"foreignKey:NumberID" json:"-"`

	CertificateNumberID *uint              `json:"certificate_number_id,omitempty"`
	CertificateNumber   *CertificateNumber `gorm:"foreignKey:CertificateNumberID" json:"-"`

	BodyPartID *uint     `json:"body_part_id,omitempty"`
	BodyPart   *BodyPart `gorm:"foreignKey:BodyPartID" json:"-"`

	MethodID *uint   `json:"method_id,omitempty"`
	Method   *Method `gorm:"foreignKey:MethodID" json:"-"`

	PlaceID *uint  `json:"place_id,omitempty"`
	Place   *Place `gorm:"foreignKey:PlaceID" json:"-"`
}

type Title struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null;uniqueIndex" json:"value"`
}

type Medication struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null;uniqueIndex" json:"value"`
}

type Dose struct {
	ID    uint    `gorm:"primarykey" json:"id"`
	Value float64 `gorm:"not null" json:"value"`
}

type Number struct {
	ID    uint `gorm:"primarykey" json:"id"`
	Value int  `gorm:"not null" json:"value"`
}

type CertificateNumber struct {
	ID    uint `gorm:"primarykey" json:"id"`
	Value int  `gorm:"not null;uniqueIndex" json:"value"`
}

type BodyPart struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null" json:"namvaluee"`
}

type Method struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null;" json:"value"`
}

type Place struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Value string `gorm:"not null;" json:"value"`
}
