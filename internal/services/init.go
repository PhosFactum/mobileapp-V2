package services

import (
	"log"

	"github.com/AlexanderMorozov1919/mobileapp/internal/config"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type Service struct {
	interfaces.ParamsParserService
	interfaces.FilterBuilderService
	interfaces.TxManager
	interfaces.ImageService
}

func NewService(db *gorm.DB, s3Cfg config.S3Config) interfaces.Service {
	parser := NewParamsParser()

	imageSvc, err := NewImageService(s3Cfg.Region, s3Cfg.BucketName, s3Cfg.Endpoint)
	if err != nil {
		log.Fatalf("Failed to initialize ImageService: %v", err)
	}

	return Service{
		ParamsParserService:  parser,
		FilterBuilderService: NewFilterBuilder(parser),
		TxManager:            NewTxManager(db),
		ImageService:         imageSvc,
	}
}
