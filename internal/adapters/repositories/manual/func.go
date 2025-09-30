package manual

import (
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
