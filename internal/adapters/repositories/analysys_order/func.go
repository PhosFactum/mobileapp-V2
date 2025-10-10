package analysis

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
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

// GetByID возвращает заказ по ID (без предзагрузки items)
func (r *AnalysisOrderRepositoryImpl) GetByID(ctx context.Context, id uint) (*entities.AnalysisOrder, error) {
	var order entities.AnalysisOrder
	err := r.GetDB(ctx).First(&order, id).Error
	return &order, err
}

// GetOrderItemsByOrderID возвращает все item'ы заказа
func (r *AnalysisOrderRepositoryImpl) GetOrderItemsByOrderID(ctx context.Context, orderID uint) ([]entities.AnalysisOrderItem, error) {
	var items []entities.AnalysisOrderItem
	err := r.GetDB(ctx).Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

// CreateOrderItems создаёт item'ы как есть (без модификации!)
func (r *AnalysisOrderRepositoryImpl) CreateOrderItems(ctx context.Context, items []entities.AnalysisOrderItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.GetDB(ctx).CreateInBatches(items, 100).Error
}

// UpdateOrderItem обновляет указанные поля item'а
func (r *AnalysisOrderRepositoryImpl) UpdateOrderItem(ctx context.Context, item entities.AnalysisOrderItem) error {
	return r.GetDB(ctx).Model(&entities.AnalysisOrderItem{}).
		Where("id = ?", item.ID).
		Select("is_completed", "completed_at", "updated_at").
		Updates(item).Error
}

// DeleteOrderItems удаляет item'ы по ID
func (r *AnalysisOrderRepositoryImpl) DeleteOrderItems(ctx context.Context, itemIDs []uint) error {
	if len(itemIDs) == 0 {
		return nil
	}
	return r.GetDB(ctx).Where("id IN ?", itemIDs).Delete(&entities.AnalysisOrderItem{}).Error
}
