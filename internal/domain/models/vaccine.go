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

	ResultID            uint `json:"result_id" binding:"required"`
	MedicationID        uint `json:"medication_id" binding:"required"`
	DoseID              uint `json:"dose_id" binding:"required"`
	NumberID            uint `json:"number_id" binding:"required"`
	CertificateNumberID uint `json:"certificate_number_id" binding:"required"`
	BodyPartID          uint `json:"body_part_id" binding:"required"`
	MethodID            uint `json:"method_id" binding:"required"`
	PlaceID             uint `json:"place_id" binding:"required"`
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
