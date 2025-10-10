package usecases

import (
	"context"
	"encoding/json"
	"fmt"

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

func (u *ReceptionUsecase) UpdateReceptionData(ctx context.Context, req *models.UpdateReceptionDataRequest) *errors.AppError {
	op := "usecase.Reception.UpdateReceptionData"

	template, err := u.repo.GetTemplateByReceptionID(ctx, req.ID)
	if err != nil {
		return errors.NewDBError(op, err)
	}

	// Шаг 2: Валидировать данные против этой схемы
	if err := ValidateJSONSchema(req.Data, template.Schema); err != nil {
		return errors.NewValidationError(op, "invalid data against template schema: "+err.Error())
	}

	// Шаг 3: Обновить данные
	if err := u.repo.UpdateReceptionData(ctx, req.ID, req.Data); err != nil {
		return errors.NewDBError(op, err)
	}

	return nil
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
