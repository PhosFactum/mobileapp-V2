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
	ExaminationType *ExaminationTypeResponse   `json:"examination_type,omitempty"`
	ExaminationView *ExaminationViewResponse   `json:"examination_view,omitempty"`
	HarmPoint       *HarmPointResponse         `json:"harm_point,omitempty"`
	PersonalInfo    *PersonalInfoResponse      `json:"personal_info,omitempty"`
	ContactInfo     *ContactInfoResponse       `json:"contact_info,omitempty"`
	Flg             *FlgResponse               `json:"flg,omitempty"`
	AnalysisOrder   *AnalysisOrderResponse     `json:"analysis_order,omitempty"`
	Statistics      *PatientStatisticsResponse `json:"statistics,omitempty"`

	Vaccines        []VaccineResponse        `json:"vaccines,omitempty"`
	Receptions      []ReceptionResponse      `json:"receptions,omitempty"`
	Specializations []SpecializationResponse `json:"specializations,omitempty"`
}

// Минимальные вложенные типы (расширьте по необходимости)
type ExaminationTypeResponse struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}

type ExaminationViewResponse struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}

type HarmPointResponse struct {
	ID    uint    `json:"id"`
	Value float32 `json:"value"`
}

type PersonalInfoResponse struct {
	ID        uint   `json:"id"`
	DocNumber string `json:"doc_number"`
	DocSeries string `json:"doc_series"`
	SNILS     string `json:"snils"`
	OMS       string `json:"oms"`

	DocumentType *DocumentTypeResponse `json:"document_type,omitempty"`
}

type DocumentTypeResponse struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
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

type CreatePatientData struct {
	FullName  string    `json:"full_name" binding:"required"`
	BirthDate time.Time `json:"birth_date" binding:"required"`
	IsMale    bool      `json:"is_male" binding:"required"`
	Position  string    `json:"position" binding:"required"`
	Division  string    `json:"division" binding:"required"`

	// Обязательные связи
	ExaminationTypeID uint `json:"examination_type_id" binding:"required"`
	ExaminationViewID uint `json:"examination_view_id" binding:"required"`
	HarmPointID       uint `json:"harm_point_id" binding:"required"`

	// Вложенные структуры
	ContactInfo  CreateContactInfoData  `json:"contact_info" binding:"required"`
	PersonalInfo CreatePersonalInfoData `json:"personal_info" binding:"required"`
}

type CreateContactInfoData struct {
	Phone   string `json:"phone" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Address string `json:"address" binding:"required"`
}

type CreatePersonalInfoData struct {
	DocNumber      string `json:"doc_number" binding:"required"`
	DocSeries      string `json:"doc_series" binding:"required"`
	SNILS          string `json:"snils" binding:"required"`
	OMS            string `json:"oms" binding:"required"`
	DocumentTypeID uint   `json:"document_type_id,omitempty"`
}
