package analysis

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

// GetPatientAnalysisOrderByID - получение конкретного направления
func (r *AnalysisRepository) GetPatientAnalysisOrderByID(patientID, orderID uint) (*entities.AnalysisOrder, error) {
	op := "repo.Analysis.GetPatientAnalysisOrderByID"

	var order entities.AnalysisOrder
	err := r.db.
		Where("patient_id = ? AND id = ?", patientID, orderID).
		First(&order).Error

	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return &order, nil
}

// GetAnalysisOrderItemsWithPagination - получение айтемов направления с пагинацией
func (r *AnalysisRepository) GetAnalysisOrderItemsWithPagination(
	orderID uint,
	page, pageSize int,
) ([]entities.AnalysisOrderItem, int64, error) {
	op := "repo.Analysis.GetAnalysisOrderItemsWithPagination"

	query := r.db.Model(&entities.AnalysisOrderItem{}).
		Preload("Analysis"). // ✅ Предзагружаем анализы для каждого айтема
		Where("order_id = ?", orderID)

	// Подсчитываем общее количество
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	// Применяем пагинацию
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	var items []entities.AnalysisOrderItem
	result := query.Order("created_at DESC").Find(&items)
	if result.Error != nil {
		return nil, 0, errors.NewDBError(op, result.Error)
	}

	return items, totalRecords, nil
}

// CreateAnalysisOrder - создание направления на анализы
func (r *AnalysisRepository) CreateAnalysisOrder(order *entities.AnalysisOrder) error {
	op := "repo.Analysis.CreateAnalysisOrder"

	if err := r.db.Create(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}

	return nil
}

// UpdateAnalysisOrder - обновление направления
func (r *AnalysisRepository) UpdateAnalysisOrder(order *entities.AnalysisOrder) error {
	op := "repo.Analysis.UpdateAnalysisOrder"

	if err := r.db.Save(order).Error; err != nil {
		return errors.NewDBError(op, err)
	}

	return nil
}

// UpdateAnalysisOrderItem - обновление элемента направления
func (r *AnalysisRepository) UpdateAnalysisOrderItem(item *entities.AnalysisOrderItem) error {
	op := "repo.Analysis.UpdateAnalysisOrderItem"

	if err := r.db.Save(item).Error; err != nil {
		return errors.NewDBError(op, err)
	}

	return nil
}

// CreateAnalysisOrderItem - создание элемента направления
func (r *AnalysisRepository) CreateAnalysisOrderItem(item *entities.AnalysisOrderItem) error {
	op := "repo.Analysis.CreateAnalysisOrderItem"

	if err := r.db.Create(item).Error; err != nil {
		return errors.NewDBError(op, err)
	}

	return nil
}

// DeleteAnalysisOrderItem - удаление элемента направления
func (r *AnalysisRepository) DeleteAnalysisOrderItem(itemID uint) error {
	op := "repo.Analysis.DeleteAnalysisOrderItem"

	result := r.db.Delete(&entities.AnalysisOrderItem{}, itemID)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}

	if result.RowsAffected == 0 {

		return errors.NewNotFoundError("analysis order item not found")
	}

	return nil
}

// UpdateAnalysisOrderTotalAmount - обновление общей суммы направления
func (r *AnalysisRepository) UpdateAnalysisOrderTotalAmount(orderID uint) error {
	op := "repo.Analysis.UpdateAnalysisOrderTotalAmount"

	// Пересчитываем общую сумму
	var totalAmount uint
	if err := r.db.Model(&entities.AnalysisOrderItem{}).
		Where("order_id = ?", orderID).
		Select("SUM(price_at_assignment)").
		Scan(&totalAmount).Error; err != nil {
		return errors.NewDBError(op, err)
	}

	// Обновляем сумму в направлении
	if err := r.db.Model(&entities.AnalysisOrder{}).
		Where("id = ?", orderID).
		Update("total_amount", totalAmount).Error; err != nil {
		return errors.NewDBError(op, err)
	}

	return nil
}
