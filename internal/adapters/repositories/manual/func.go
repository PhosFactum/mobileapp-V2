package manual

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *ManualRepositoryImpl) GetManualValueByTypeAndID(id uint, ref_type entities.ReferenceType) (string, error) {
	op := "repo.Manual.GetDocumentTypeValueByID"
	if id == 0 {
		return "", nil
	}
	var entry entities.Manual
	err := r.db.
		Where("id = ? AND type = ?", id, ref_type).
		First(&entry).
		Error
	if err != nil {
		return "", errors.NewDBError(op, err)
	}
	return entry.Value, nil
}

func (r *ManualRepositoryImpl) GetAllManuals(ctx context.Context) ([]entities.Manual, error) {
	var manuals []entities.Manual
	if err := r.db.WithContext(ctx).Find(&manuals).Error; err != nil {
		return nil, errors.NewDBError("repo.Manual.GetAllManuals", err)
	}
	return manuals, nil
}
