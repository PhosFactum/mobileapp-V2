package flg

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

func (r *FLGRepositoryImpl) CreateFLG(flg entities.FLG) (uint, error) {
	if err := r.db.Create(&flg).Error; err != nil {
		return 0, errors.NewDBError("repo.FLG.CreateFLG", err)
	}
	return flg.ID, nil
}

func (r *FLGRepositoryImpl) UpdateFLG(id uint, updateMap map[string]interface{}) (*entities.FLG, error) {
	var flg entities.FLG
	if err := r.db.First(&flg, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.NewDBError("repo.FLG.UpdateFLG", err)
	}

	if err := r.db.Model(&flg).Updates(updateMap).Error; err != nil {
		return nil, errors.NewDBError("repo.FLG.UpdateFLG", err)
	}

	return &flg, nil
}
