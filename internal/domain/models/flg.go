package models

import "mime/multipart"

type FlgResponse struct {
	ID           uint   `json:"id"`
	Organization string `json:"organization" example:"Stavropol"`
	Number       string `json:"number" example:"984212"`
	Result       string `json:"result" example:"COVID"`
	Date         string `json:"date" example:"2023-10-15T14:30:00Z"`
	PhotoURL     string `json:"photo_url"` // ← временный URL для просмотра
}

// CreateFlgWithPhotoRequest — multipart-запрос
type CreateFlgWithPhotoRequest struct {
	PatientID    uint                  `form:"patient_id" binding:"required"`
	Organization string                `form:"organization" binding:"required"`
	Number       string                `form:"number" binding:"required"`
	Result       string                `form:"result" binding:"required"`
	Date         string                `form:"date" binding:"required,datetime=2006-01-02"`
	File         *multipart.FileHeader `form:"file" binding:"required"`
}
