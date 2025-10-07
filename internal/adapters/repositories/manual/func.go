package manual

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

// internal/adapters/repositories/manual/manual_repository.go
func (r *ManualRepositoryImpl) GetManualValueByTypeAndID(ctx context.Context, id uint, refType entities.ReferenceType) (string, error) {
	op := "repo.Manual.GetManualValueByTypeAndID"
	if id == 0 {
		return "", nil
	}
	var entry entities.Manual
	err := r.GetDB(ctx).WithContext(ctx).
		Where("id = ? AND type = ?", id, refType).
		First(&entry).
		Error
	if err != nil {
		return "", errors.NewDBError(op, err)
	}
	return entry.Value, nil
}

func (r *ManualRepositoryImpl) GetAllManuals(ctx context.Context) ([]entities.Manual, error) {
	op := "repo.Manual.GetAllManuals"
	var manuals []entities.Manual
	if err := r.GetDB(ctx).WithContext(ctx).Find(&manuals).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return manuals, nil
}

func (r *ManualRepositoryImpl) GetManualValuesByType(ctx context.Context, refType entities.ReferenceType) ([]string, error) {
	op := "repo.Manual.GetManualValuesByType"
	var values []string
	err := r.GetDB(ctx).WithContext(ctx).
		Model(&entities.Manual{}).
		Where("type = ?", refType).
		Pluck("value", &values).
		Error
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return values, nil
}
