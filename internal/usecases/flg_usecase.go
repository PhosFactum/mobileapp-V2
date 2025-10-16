package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/google/uuid"
)

type FlgUseсase struct {
	repo         interfaces.FlgRepository
	imageService interfaces.ImageService
}

func NewFlgUseсase(repo interfaces.FlgRepository) interfaces.FlgUsecase {
	return &FlgUseсase{repo: repo}
}

func (u *FlgUseсase) CreateFlgWithPhoto(ctx context.Context, req *models.CreateFlgRequest) (*models.FlgResponse, *errors.AppError) {
	op := "usecase.Flg.CreateFlgWithPhoto"

	// 1. Проверка пациента
	flgs, err := u.repo.GetFlgByPatientID(ctx, req.PatientID)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	if len(flgs) == 0 {
		return nil, errors.NewNotFoundError(op)
	}

	// 2. Валидация Content-Type
	if !isValidImageContentType(req.ContentType) {
		return nil, errors.NewValidationError(op, "invalid content type: only JPEG/PNG allowed")
	}

	// 3. Ограничение размера (10 МБ)
	const maxFileSize = 10 << 20
	if len(req.FileData) > maxFileSize {
		return nil, errors.NewValidationError(op, "file too large (max 10 MB)")
	}

	// 4. Парсинг даты
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.NewValidationError(op, "invalid date format, expected YYYY-MM-DD")
	}

	// 5. Генерация ключа в S3
	ext := getExtensionFromContentType(req.ContentType)
	key := fmt.Sprintf("flg/%d/%s%s", req.PatientID, uuid.NewString(), ext)

	// 6. Загрузка в S3
	if err := u.imageService.UploadObject(ctx, key, req.ContentType, req.FileData); err != nil {
		return nil, errors.NewInternalError(op, "failed to upload image to storage", err)
	}

	// 7. Сохранение в БД
	flg := &entities.Flg{
		PatientID:    req.PatientID,
		Organization: req.Organization,
		Number:       req.Number,
		Result:       req.Result,
		Date:         parsedDate,
		PhotoKey:     key,
	}

	if err := u.repo.CreateFlg(ctx, flg); err != nil {
		// Откат: удаляем из S3
		_ = u.imageService.DeleteObject(ctx, key)
		return nil, errors.NewDBError(op, err)
	}

	// 8. Генерация временного URL
	presignedURL, err := u.imageService.GetPresignedURL(ctx, key)
	if err != nil {
		fmt.Printf("Warning: failed to generate presigned URL: %v\n", err)
		presignedURL = ""
	}

	return &models.FlgResponse{
		ID:           flg.ID,
		Organization: flg.Organization,
		Number:       flg.Number,
		Result:       flg.Result,
		Date:         flg.Date.Format("2006-01-02"),
		PhotoURL:     presignedURL,
	}, nil
}

// Вспомогательные функции
func isValidImageContentType(ct string) bool {
	ct = strings.ToLower(ct)
	return ct == "image/jpeg" || ct == "image/jpg" || ct == "image/png"
}

func getExtensionFromContentType(ct string) string {
	switch strings.ToLower(ct) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ".jpg"
	}
}
