package entities

import "time"

type Vaccine struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	TitleID     uint      `json:"title_id,omitempty"`
	PatientID   uint      `gorm:"index"`

	ResultID            uint `json:"result"`
	MedicationID        uint `json:"medication_id,omitempty"`
	DoseID              uint `json:"dose_id,omitempty"`
	NumberID            uint `json:"number_id,omitempty"`
	CertificateNumberID uint `json:"certificate_number_id,omitempty"`
	BodyPartID          uint `json:"body_part_id,omitempty"`
	MethodID            uint `json:"method_id,omitempty"`
	PlaceID             uint `json:"place_id,omitempty"`
}

type VaccineRefusal struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	Date        time.Time `json:"date"`
	TitleID     uint      `json:"title_id,omitempty"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	PatientID   uint      `gorm:"index"`
}

type VaccineWithdrawal struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	Date        time.Time `json:"date"`
	TitleID     uint      `json:"title_id,omitempty"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	PatientID   uint      `gorm:"index"`
	Num         int       `json:"med_withdrawl_num"`
}

type Titr struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	Date        time.Time `json:"date"`
	TitleID     uint      `json:"title_id,omitempty"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	PatientID   uint      `gorm:"index"`
	Amount      string    `json:"titer_amount"`
}
