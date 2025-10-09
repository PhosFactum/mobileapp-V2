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
