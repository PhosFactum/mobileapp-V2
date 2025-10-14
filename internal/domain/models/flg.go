package models

type FlgResponse struct {
	ID           uint   `json:"id"`
	IsCompleted  bool   `json:"is_completed"`
	Organization string `json:"organization" example:"Stavropol"`
	Number       string `json:"number" example:"984212"`
	Result       string `json:"result" example:"COVID"`
	Date         string `json:"date" example:"2023-10-15T14:30:00Z"`
}
