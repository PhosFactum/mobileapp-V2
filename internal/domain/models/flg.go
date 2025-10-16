package models

// CreateFlgRequest — данные для создания FLG (уже без multipart!)
type CreateFlgRequest struct {
	PatientID    uint   `json:"patient_id"`
	Organization string `json:"organization"`
	Number       string `json:"number"`
	Result       string `json:"result"`
	Date         string `json:"date"` // "2025-10-14"

	// Данные изображения
	FileData    []byte `json:"-"` // не сериализуется в JSON
	ContentType string `json:"-"` // например: "image/jpeg"
}

// FlgResponse — ответ
type FlgResponse struct {
	ID           uint   `json:"id"`
	Organization string `json:"organization"`
	Number       string `json:"number"`
	Result       string `json:"result"`
	Date         string `json:"date"`
	PhotoURL     string `json:"photo_url"` // временный URL
}
