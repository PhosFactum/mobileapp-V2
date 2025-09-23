package analysis

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

// GetAnalysisOrderItemsByOrderID возвращает все элементы направления на анализы по order_id с предзагрузкой анализа
func (r *AnalysisRepositoryImpl) GetAnalysisOrderItemsByOrderID(orderID uint) ([]entities.AnalysisOrderItem, error) {
	op := "repo.AnalysisOrder.GetAnalysisOrderItemsByOrderID"

	var items []entities.AnalysisOrderItem

	err := r.db.
		Where("order_id = ?", orderID).
		Preload("Analysis").
		Order("created_at DESC").
		Find(&items).Error

	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return items, nil
}
