package models

type PatientGroupShortResponse struct {
	ID                uint   `json:"id"`
	CreatedAt         string `json:"created_at" example:"2023-05-15T14:30:00Z"`
	Code              string `json:"code" example:"94928490"`
	OrganizationTitle string `json:"organization" example:"medlab"`
}
