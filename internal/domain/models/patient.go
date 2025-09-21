package models

import (
	"time"
)

// UpdatePatientRequest - запрос на обновление данных пациента
// @Description Данные для обновления информации о пациенте
type UpdatePatientRequest struct {
	ID        uint   `json:"id" example:"10"`                              // ID пациента
	FullName  string `json:"full_name" example:"Смирнов Алексей Петрович"` // ФИО пациента
	BirthDate string `json:"birth_date" example:"1980-05-15"`              // Дата рождения
	IsMale    bool   `json:"is_male" example:"true"`                       // Пол (true - мужской)
}

// CreatePatientRequest - запрос на создание пациента
// @Description Данные для создания нового пациента
type CreatePatientRequest struct {
	FullName  string `json:"full_name" example:"Смирнов Алексей Петрович"` // ФИО пациента
	BirthDate string `json:"birth_date" example:"1980-05-15"`              // Дата рождения
	IsMale    bool   `json:"is_male" example:"true"`                       // Пол (true - мужской)
}

// ShortPatientResponse - краткая информация о пациенте
// @Description Сокращенные данные пациента
type ShortPatientResponse struct {
	ID        uint      `json:"id" example:"1"`
	FullName  string    `json:"full_name" example:"Смирнов Алексей Петрович"`
	BirthDate time.Time `json:"birth_date" example:"1980-05-15T00:00:00Z"` // Дата рождения
	IsMale    bool      `json:"is_male" example:"true"`                    // Пол (true - мужской)
}

// PatientResponse - полная информация о пациенте
// @Description Все данные пациента
type PatientResponse struct {
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
	PatientGroupID    uint `json:"patient_group_id" binding:"required"`

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
