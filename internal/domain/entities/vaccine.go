package entities

import "time"

type Vaccine struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	PatientID   uint      `gorm:"index"`

	// Полная прививка
	Result            string `json:"result"`
	Title             string `json:"title_id,omitempty"`
	Medication        string `json:"medication_id,omitempty"`
	Dose              string `json:"dose_id,omitempty"`
	Number            string `json:"number_id,omitempty"`
	CertificateNumber string `json:"certificate_number_id,omitempty"`
	BodyPart          string `json:"body_part_id,omitempty"`
	Method            string `json:"method_id,omitempty"`
	Place             string `json:"place_id,omitempty"`
}

type VaccineRefusal struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	PatientID   uint      `gorm:"index"`
}

type VaccineWithdrawal struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	PatientID   uint      `gorm:"index"`
	Num         int       `json:"med_withdrawl_num"`
}

type Titr struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	PatientID   uint      `gorm:"index"`
	Amount      int       `json:"titer_amount"`
}
