package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
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

func (u *ReceptionUsecase) CreateReception(ctx context.Context, req *models.CreateReceptionRequest) (*entities.Reception, *errors.AppError) {

	op := "usecase.Reception.CreateReception"

	// 1. Получить актуальную схему с сервера (не доверяем клиенту!)
	schema, err := u.repo.GetTemplateSchemaByID(ctx, req.TemplateID)
	if err != nil {
		return nil, errors.NewDBError(op, err) // AppError
	}

	// 2. Валидировать данные
	if err := ValidateJSONSchema(req.Data, schema); err != nil {
		return nil, errors.NewValidationError(op, err.Error())
	}

	// 3. Создать приём
	reception := &entities.Reception{
		TemplateID:       req.TemplateID,
		Data:             req.Data,
		PatientID:        req.PatientID,
		SpecializationID: req.SpecializationID,
		// CreatedAt, UpdatedAt — GORM заполнит автоматически
	}

	if err := u.repo.CreateReception(ctx, reception); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return reception, nil
}

// вынести
func ValidateJSONSchema(data, schema []byte) error {
	compiled, err := jsonschema.CompileString("reception-schema", string(schema))
	if err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON  %w", err)
	}

	if err := compiled.Validate(jsonData); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
