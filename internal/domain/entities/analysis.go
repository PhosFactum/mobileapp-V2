package entities

import "time"

type Analysis struct {
	ID uint `gorm:"primarykey" json:"id" example:"1"`

	Code  string `gorm:"not null" json:"code" example:"12-0739"`
	Title string `gorm:"not null" json:"title" example:"EKG"`
	Price uint   `gorm:"not null" json:"price" example:"100"`
}

type AnalysisOrderItem struct {
	ID uint `gorm:"primarykey" json:"id"`

	OrderID uint           `gorm:"not null;index" json:"order_id"`
	Order   *AnalysisOrder `gorm:"foreignKey:OrderID" json:"-"`

	AnalysisID uint      `gorm:"not null;index" json:"analysis_id"`
	Analysis   *Analysis `gorm:"foreignKey:AnalysisID" json:"analysis"`

	// Статус конкретного анализа
	IsCompleted bool       `gorm:"default:false" json:"is_completed"` // Сдан или нет
	CompletedAt *time.Time `json:"completed_at,omitempty"`            // Когда сдан

	PriceAtAssignment uint `gorm:"not null" json:"price_at_assignment"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AnalysisOrder - направление на анализы (промежуточная структура)
type AnalysisOrder struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	OrderNumber string    `gorm:"not null;uniqueIndex" json:"order_number"` // Номер направления (уникальный)

	PatientID uint `gorm:"not null;index" json:"patient_id"`

	OrderItems []AnalysisOrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}
