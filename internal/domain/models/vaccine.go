package models

import (
	"time"
)

// Возможно добавить количество титров и прочее
type VaccineAllResponse struct {
	ID             uint      `json:"id"`
	Date           time.Time `json:"date"`
	Type           string    `json:"type"`
	Title          string    `json:"title"`
	TiterAmountStr *string   `json:"titer_amount_str,omitempty"`
}
