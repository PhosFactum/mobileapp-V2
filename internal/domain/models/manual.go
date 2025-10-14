package models

type ManualResponse struct {
	ID    uint   `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
