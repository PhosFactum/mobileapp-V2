package services

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ImageService struct {
	client     *s3.Client
	bucketName string
}

// NewImageService создаёт S3-клиент
func NewImageService(region, bucketName, endpoint string) (*ImageService, error) {
	// 1. Загружаем базовую конфигурацию (credentials, region и т.д.)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// 2. Создаём S3-клиент с endpoint (если указан)
	var s3Client *s3.Client
	if endpoint != "" {
		// Используем endpoint напрямую в S3-клиенте
		s3Client = s3.NewFromConfig(cfg, s3.WithEndpointResolver(
			s3.EndpointResolverFunc(func(region string, opts s3.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: endpoint,
				}, nil
			}),
		))
	} else {
		// Стандартный клиент для AWS
		s3Client = s3.NewFromConfig(cfg)
	}

	return &ImageService{
		client:     s3Client,
		bucketName: bucketName,
	}, nil
}

// UploadObject загружает БАЙТЫ в S3
func (s *ImageService) UploadObject(ctx context.Context, key, contentType string, data []byte) error {
	_, err := manager.NewUploader(s.client).Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	return err
}

// GetPresignedURL генерирует временный URL
func (s *ImageService) GetPresignedURL(ctx context.Context, key string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return req.URL, nil
}
