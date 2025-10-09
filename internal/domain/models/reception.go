package models

import "encoding/json"

type ReceptionResponse struct {
	ID               uint                      `json:"id"`
	IsCompleted      bool                      `json:"is_completed"`
	SpecializationID uint                      `json:"specialization_id"`
	Specialization   SpecializationResponse    `json:"specialization,omitempty"`
	Template         ReceptionTemplateResponse `json:"template"`
	Data             json.RawMessage           `json:"data"`
}

type ReceptionTemplateResponse struct {
	ID     uint            `json:"id"`
	Code   string          `json:"code"`
	Schema json.RawMessage `json:"fields"`
}

type CreateReceptionRequest struct {
	IsCompleted      bool            `json:"is_completed" binding:"required"`
	TemplateID       uint            `json:"template_id" binding:"required"`
	Data             json.RawMessage `json:"data" binding:"required"`
	PatientID        uint            `json:"patient_id" binding:"required"`
	SpecializationID uint            `json:"specialization_id" binding:"required"`
}
