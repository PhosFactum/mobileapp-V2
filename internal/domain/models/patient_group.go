package models

type PatientGroupShortResponse struct {
	ID                uint   `json:"id"`
	CreatedAt         string `json:"created_at" example:"2023-05-15T14:30:00Z"`
	Code              string `json:"code" example:"94928490"`
	OrganizationTitle string `json:"organization" example:"medlab"`
}

// PatientGroupWithPatientsResponse - ответ с группой пациентов и списком пациентов
type PatientGroupWithPatientsResponse struct {
	ID           uint                   `json:"id"`
	CreatedAt    string                 `json:"created_at"`
	Code         string                 `json:"code"`
	Organization string                 `json:"organization"`
	Patients     []ShortPatientResponse `json:"patients"`
}

// OrganizationWithPatientGroupsResponse - ответ с организацией и группами пациентов
type OrganizationWithPatientGroupsResponse struct {
	OrganizationID   uint                               `json:"organization_id"`
	OrganizationName string                             `json:"organization_name"`
	PatientGroups    []PatientGroupWithPatientsResponse `json:"patient_groups"`
}
