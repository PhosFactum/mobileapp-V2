package models

import "encoding/json"

type ReceptionResponse struct {
	ID               uint                    `json:"id"`
	IsCompleted      bool                    `json:"is_completed"`
	SpecializationID uint                    `json:"specialization_id"`
	Specialization   *SpecializationResponse `json:"specialization,omitempty"`
	Data             json.RawMessage         `json:"data"`
}
