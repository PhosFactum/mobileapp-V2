package models

import (
	"time"
)

type PatientResponse struct {
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	BirthDate time.Time `json:"birth_date"`
	Age       int       `json:"age"`
	IsMale    bool      `json:"is_male"`
	Position  string    `json:"position"`
	Division  string    `json:"division"`

	PatientGroupID uint `json:"patient_group_id"`

	// Вложенные объекты
	ExaminationTypeID uint `json:"examination_type,omitempty"`
	ExaminationViewID uint `json:"examination_view,omitempty"`

	HarmPoint     HarmPointResponse          `json:"harm_point"`
	PersonalInfo  PersonalInfoResponse       `json:"personal_info"`
	ContactInfo   ContactInfoResponse        `json:"contact_info"`
	AnalysisOrder AnalysisOrderResponse      `json:"analysis_order"`
	Statistics    *PatientStatisticsResponse `json:"statistics,omitempty"`
	Flg           *FlgResponse               `json:"flg,omitempty"`

	// Связанные коллекции
	Vaccines        []VaccineAllResponse     `json:"vaccines,omitempty"`
	Receptions      []ReceptionResponse      `json:"receptions,omitempty"`
	Specializations []SpecializationResponse `json:"specializations,omitempty"`
}

type HarmPointResponse struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}

type PersonalInfoResponse struct {
	ID        uint   `json:"id"`
	DocNumber string `json:"doc_number"`
	DocSeries string `json:"doc_series"`
	SNILS     string `json:"snils"`
	OMS       string `json:"oms"`

	DocumentTypeID uint `json:"document_type,omitempty"`
}

type ContactInfoResponse struct {
	ID      uint   `json:"id"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type PatientStatisticsResponse struct {
	ID                     uint  `json:"id"`
	TotalReceptions        int64 `json:"total_receptions"`
	CompletedReceptions    int64 `json:"completed_receptions"`
	TotalAnalysisOrders    int64 `json:"total_analysis_orders"`
	CompletedAnalysisItems int64 `json:"completed_analysis_items"`
}

type SpecializationResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type CreatePatientRequest struct {
	FullName          string    `json:"full_name" binding:"required"`
	BirthDate         time.Time `json:"birth_date" binding:"required"`
	IsMale            bool      `json:"is_male" binding:"required"`
	Position          string    `json:"position" binding:"required"`
	Division          string    `json:"division" binding:"required"`
	ExaminationTypeID uint      `json:"examination_type_id" binding:"required"`
	ExaminationViewID uint      `json:"examination_view_id" binding:"required"`
	GroupID           uint      `json:"group_id" binding:"required"`

	// Обязательные связи
	HarmPointID uint `json:"harm_point_id" binding:"required"`

	// Вложенные структуры
	ContactInfo  CreateContactInfoRequest  `json:"contact_info" binding:"required"`
	PersonalInfo CreatePersonalInfoRequest `json:"personal_info" binding:"required"`
}

type CreateContactInfoRequest struct {
	Phone   string `json:"phone" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Address string `json:"address" binding:"required"`
}

type CreatePersonalInfoRequest struct {
	DocNumber      string `json:"doc_number" binding:"required"`
	DocSeries      string `json:"doc_series" binding:"required"`
	SNILS          string `json:"snils" binding:"required"`
	OMS            string `json:"oms" binding:"required"`
	DocumentTypeID uint   `json:"document_type_id,omitempty"`
}
