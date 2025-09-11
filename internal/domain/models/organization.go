package models

import "github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"

type OrganizationShortResponse struct {
	ID      uint             `json:"id"`
	Title   string           `json:"title"`
	Manager entities.Manager `json:"manager"`
}
