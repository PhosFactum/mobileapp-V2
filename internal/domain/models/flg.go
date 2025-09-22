package models

import "time"

type FLGCreateRequest struct {
	PatientID       uint      `json:"patient_id" binding:"required"`
	OrganizationID  uint      `json:"organization_id" binding:"required"`
	ExaminationDate time.Time `json:"examination_date" binding:"required"`
	Number          string    `json:"number" binding:"required"`
	Result          string    `json:"result" binding:"required"`
	AttachedImage   string    `json:"attached_image,omitempty"`
}

type FLGUpdateRequest struct {
	OrganizationID  uint      `json:"organization_id,omitempty"`
	ExaminationDate time.Time `json:"examination_date,omitempty"`
	Number          string    `json:"number,omitempty"`
	Result          string    `json:"result,omitempty"`
	AttachedImage   string    `json:"attached_image,omitempty"`
}

type FLGResponse struct {
	ID              uint      `json:"id"`
	PatientID       uint      `json:"patient_id"`
	OrganizationID  uint      `json:"organization_id"`
	ExaminationDate time.Time `json:"examination_date"`
	Number          string    `json:"number"`
	Result          string    `json:"result"`
	AttachedImage   string    `json:"attached_image,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
