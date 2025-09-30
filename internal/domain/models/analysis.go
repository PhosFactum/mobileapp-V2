package models

type AnalysisOrderResponse struct {
	ID          uint                        `json:"id"`
	OrderNumber string                      `json:"order_number"`
	TotalAmount uint                        `json:"total_amount"`
	OrderItems  []AnalysisOrderItemResponse `json:"order_items"`
}

type AnalysisOrderItemResponse struct {
	ID          uint              `json:"id"`
	AnalysisID  uint              `json:"analysis_id"`
	Analysis    *AnalysisResponse `json:"analysis"`
	IsCompleted bool              `json:"is_completed"`
}

type AnalysisResponse struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Title string `json:"title"`
	Price uint   `json:"price"`
}
