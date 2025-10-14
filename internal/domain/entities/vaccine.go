package entities

import "time"

type Vaccine struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	Date      time.Time `json:"date"`
	TitleID   uint      `gorm:"not null" json:"title_id,omitempty"`
	PatientID uint      `gorm:"index"`

	ResultID            uint `json:"result"`
	MedicationID        uint `json:"medication_id,omitempty"`
	DoseID              uint `json:"dose_id,omitempty"`
	NumberID            uint `json:"number_id,omitempty"`
	CertificateNumberID uint `json:"certificate_number_id,omitempty"`
	BodyPartID          uint `json:"body_part_id,omitempty"`
	MethodID            uint `json:"method_id,omitempty"`
	PlaceID             uint `json:"place_id,omitempty"`

	PhotoURL *string `json:"photo_url,omitempty" example:"https://my-bucket.s3.amazonaws.com/flg/123/photo.jpg"`
}

type VaccineRefusal struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	Date      time.Time `json:"date"`
	TitleID   uint      `gorm:"not null" json:"title_id"`
	PatientID uint      `gorm:"index"`

	PhotoURL *string `json:"photo_url,omitempty" example:"https://my-bucket.s3.amazonaws.com/flg/123/photo.jpg"`
}

type VaccineWithdrawal struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	Date      time.Time `json:"date"`
	TitleID   uint      `gorm:"not null" json:"title_id"`
	PatientID uint      `gorm:"index"`
	Num       int       `json:"med_withdrawl_num"`

	PhotoURL *string `json:"photo_url,omitempty" example:"https://my-bucket.s3.amazonaws.com/flg/123/photo.jpg"`
}

type Titr struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	Date      time.Time `json:"date"`
	TitleID   uint      `gorm:"not null" json:"title_id"`
	PatientID uint      `gorm:"index"`
	Amount    string    `json:"titer_amount"`

	PhotoURL *string `json:"photo_url,omitempty" example:"https://my-bucket.s3.amazonaws.com/flg/123/photo.jpg"`
}
