package models

type DoctorResponse struct {
	ID              uint                     `json:"ID" example:"1"`                           // id
	FullName        string                   `json:"full_name" example:"Иванов Иван Иванович"` // Полное имя врача
	Specializations []SpecializationResponse `json:"specializations,omitempty"`
}

type DoctorLoginRequest struct {
	Phone    string `json:"phone" binding:"required" example:"+79622840765"` // Логин (телефон)
	Password string `json:"password" binding:"required" example:"123"`       // Пароль
}

// DoctorAuthResponse - ответ на авторизацию врача
// @Description Ответ с данными авторизованного врача
type DoctorAuthResponse struct {
	ID    uint   `json:"id" example:"1"`                // ID врача
	Token string `json:"token" example:"eyJhbGciOi..."` // JWT токен
}

type DoctorInfoResponse struct {
	DoctorID       uint   `json:"doctor_id" example:"1"`
	FullName       string `json:"full_name" example:"Иванов Иван Иванович"` // Полное имя врача
	Specialization string `json:"specialization"`
}

// CreateDoctorRequest - запрос на создание врача
// @Description Используется для регистрации нового врача в системе
type CreateDoctorRequest struct {
	FullName         string `json:"full_name" binding:"required" example:"Иванов Иван Иванович"` // ФИО врача (обязательное)
	Phone            string `json:"phone" binding:"required" example:"+79123456789"`             // Логин (обязательное)
	Password         string `json:"password" binding:"required" example:"qwerty123"`             // Пароль (обязательное)
	SpecializationID uint   `json:"specialization_id" binding:"required" example:"1"`            // Специализация (обязательное)
}

// UpdateDoctorRequest - запрос на обновление данных врача
// @Description Используется для изменения информации о враче
type UpdateDoctorRequest struct {
	ID               uint   `json:"id" example:"1"`
	FullName         string `json:"full_name" example:"Иванов Иван Иванович"`
	Phone            string `json:"phone" example:"+79123456789"`
	PasswordHash     string `json:"-"` // Убрали из JSON, чтобы не принималось извне
	SpecializationID uint   `json:"specialization_id" example:"1"`
}
