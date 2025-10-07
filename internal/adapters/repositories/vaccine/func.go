package vaccine

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

// CreateVaccine создаёт запись о вакцине
func (r *VaccineRepositoryImpl) CreateVaccine(ctx context.Context, vaccine *entities.Vaccine) error {
	op := "repo.Vaccine.CreateVaccine"
	if err := r.GetDB(ctx).WithContext(ctx).Create(vaccine).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreateVaccineRefusal создаёт запись об отказе от вакцины
func (r *VaccineRepositoryImpl) CreateVaccineRefusal(ctx context.Context, refusal *entities.VaccineRefusal) error {
	op := "repo.Vaccine.CreateVaccineRefusal"
	if err := r.GetDB(ctx).WithContext(ctx).Create(refusal).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreateVaccineWithdrawal создаёт запись об отводе от вакцины
func (r *VaccineRepositoryImpl) CreateVaccineWithdrawal(ctx context.Context, withdrawal *entities.VaccineWithdrawal) error {
	op := "repo.Vaccine.CreateVaccineWithdrawal"
	if err := r.GetDB(ctx).WithContext(ctx).Create(withdrawal).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreateTitr создаёт запись о титровании
func (r *VaccineRepositoryImpl) CreateTitr(ctx context.Context, titr *entities.Titr) error {
	op := "repo.Vaccine.CreateTitr"
	if err := r.GetDB(ctx).WithContext(ctx).Create(titr).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}
