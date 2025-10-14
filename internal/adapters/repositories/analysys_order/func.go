package analysis

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *AnalysisOrderRepositoryImpl) CreateAnalysisItems(ctx context.Context, items []entities.AnalysisOrderItem) error {
	op := "repo.Patient.CreateAnalysisItems"
	if err := r.GetDB(ctx).WithContext(ctx).Create(&items).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// CreateAnalysisOrder создаёт направление на анализы
func (r *AnalysisOrderRepositoryImpl) CreateAnalysisOrder(ctx context.Context, order *entities.AnalysisOrder) error {
	op := "repo.Patient.CreateAnalysisOrder"
	if err := r.GetDB(ctx).WithContext(ctx).Create(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

// UpdateAnalysisOrder обновляет направление
func (r *AnalysisOrderRepositoryImpl) UpdateAnalysisOrder(ctx context.Context, order *entities.AnalysisOrder) error {
	op := "repo.Patient.UpdateAnalysisOrder"
	if err := r.GetDB(ctx).WithContext(ctx).Save(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *AnalysisOrderRepositoryImpl) GetByID(ctx context.Context, id uint) (*entities.AnalysisOrder, error) {
	op := "repo.AnalysisOrder.GetByID"
	var order entities.AnalysisOrder
	if err := r.GetDB(ctx).First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError(op)
		}
		return nil, errors.NewDBError(op, err)
	}
	return &order, nil
}

func (r *AnalysisOrderRepositoryImpl) GetOrderItemsByOrderID(ctx context.Context, orderID uint) ([]entities.AnalysisOrderItem, error) {
	op := "repo.AnalysisOrder.GetOrderItemsByOrderID"
	var items []entities.AnalysisOrderItem
	if err := r.GetDB(ctx).Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return items, nil
}

func (r *AnalysisOrderRepositoryImpl) UpsertOrderItems(ctx context.Context, items []entities.AnalysisOrderItem) error {
	op := "repo.AnalysisOrder.UpsertOrderItems"
	if len(items) == 0 {
		return nil
	}

	err := r.GetDB(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "order_id"},
			{Name: "analysis_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"is_completed", "completed_at", "updated_at"}),
	}).Create(&items).Error

	if err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}
