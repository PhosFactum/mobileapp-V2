package flg

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *FlgRepositoryImpl) CreateFlg(ctx context.Context, flg *entities.Flg) error {
	op := "repo.Flg.CreateFlg"
	if err := r.GetDB(ctx).WithContext(ctx).Create(flg).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *FlgRepositoryImpl) GetFlgByPatientID(ctx context.Context, patientID uint) ([]entities.Flg, error) {
	var flgs []entities.Flg
	if err := r.GetDB(ctx).WithContext(ctx).
		Where("patient_id = ?", patientID).
		Find(&flgs).Error; err != nil {
		return nil, errors.NewDBError("repo.Flg.GetByPatientID", err)
	}
	return flgs, nil
}

func (r *FlgRepositoryImpl) Delete(ctx context.Context, id uint) error {
	op := "repo.Flg.Delete"
	if err := r.GetDB(ctx).WithContext(ctx).Delete(&entities.Flg{}, id).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}
