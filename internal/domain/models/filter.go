package models

// FilterResponse - модель для фильтрации
// @Description Выводит информацию о пагинации и фильтрации
type FilterResponse[T any] struct {
	Hits        T   `json:"hits"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalHits   int `json:"total_hits"`
	HitsPerPage int `json:"hits_per_page"`
}
