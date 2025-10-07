package models

import (
	"time"
)

// Возможно добавить количество титров и прочее
type VaccineAllResponse struct {
	ID             uint      `json:"id"`
	Date           time.Time `json:"date"`
	Type           string    `json:"type"`
	Title          string    `json:"title"`
	TiterAmountStr *string   `json:"titer_amount_str,omitempty"`
}

type CreateVaccineRequest struct {
	Date      time.Time `json:"date" binding:"required"`
	TitleID   uint      `json:"title_id" binding:"required"`
	PatientID uint      `json:"patient_id" binding:"required"`

	ResultID            uint  `json:"result_id" binding:"required"`
	MedicationID        *uint `json:"medication_id,omitempty"`
	DoseID              *uint `json:"dose_id,omitempty"`
	NumberID            *uint `json:"number_id,omitempty"`
	CertificateNumberID *uint `json:"certificate_number_id,omitempty"`
	BodyPartID          *uint `json:"body_part_id,omitempty"`
	MethodID            *uint `json:"method_id,omitempty"`
	PlaceID             *uint `json:"place_id,omitempty"`
}

type CreateVaccineRefusalRequest struct {
	Date      time.Time `json:"date" binding:"required"`
	TitleID   uint      `json:"title_id" binding:"required"`
	PatientID uint      `json:"patient_id" binding:"required"`
}

type CreateVaccineWithdrawalRequest struct {
	Date      time.Time `json:"date" binding:"required"`
	TitleID   uint      `json:"title_id" binding:"required"`
	PatientID uint      `json:"patient_id" binding:"required"`
	Num       int       `json:"med_withdrawl_num" binding:"required"`
}

type CreateTitrRequest struct {
	Date      time.Time `json:"date" binding:"required"`
	TitleID   uint      `json:"title_id" binding:"required"`
	PatientID uint      `json:"patient_id" binding:"required"`
	Amount    string    `json:"titer_amount" binding:"required"`
}
