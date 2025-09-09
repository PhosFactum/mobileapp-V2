package models

type OrganizationShortResponse struct {
	ID             uint   `json:"id"`
	Title          string `json:"title"`
	Code           string `json:"code"`
	DoctorFullName string `json:"doctor_full_name"`
}
