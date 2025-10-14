package models

import "github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"

// обёртки для Swagger-документации

type PatientsListResponse struct {
	Hits        []entities.Patient `json:"hits"`
	CurrentPage int                `json:"currentPage"`
	TotalPages  int                `json:"totalPages"`
	TotalHits   int                `json:"totalHits"`
	HitsPerPage int                `json:"hitsPerPage"`
}
