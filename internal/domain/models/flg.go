package models

type FlgResponse struct {
	ID           uint   `json:"id"`
	IsCompleted  bool   `json:"is_completed"`
	Organization string `json:"organization" example:"Stavropol"`
	Number       string `json:"number" example:"984212"`
	Result       string `json:"result" example:"COVID"`
}
