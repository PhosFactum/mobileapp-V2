package usecases

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
)

type ReceptionUsecase struct {
	repo          interfaces.ReceptionRepository
	FilterBuilder interfaces.FilterBuilderService
}

func NewReceptionUsecase(repo interfaces.ReceptionRepository, s interfaces.Service) interfaces.ReceptionUsecase {
	return &ReceptionUsecase{
		repo:          repo,
		FilterBuilder: s}
}
