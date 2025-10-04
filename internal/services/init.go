package services

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type Service struct {
	interfaces.ParamsParserService
	interfaces.FilterBuilderService
	interfaces.TxManager
}

func NewService(db *gorm.DB) interfaces.Service {
	parser := NewParamsParser()
	return Service{
		parser,
		NewFilterBuilder(parser),
		NewTxManager(db),
	}
}
