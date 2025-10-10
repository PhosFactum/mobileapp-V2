package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

type AnalysisOrderUsecase struct {
	analysisRepo      interfaces.AnalysisRepository
	analysisOrderRepo interfaces.AnalysisOrderRepository
	txManager         interfaces.TxManager
}

func NewAnalysisOrderUsecase(analysisRepo interfaces.AnalysisRepository, analysisOrderRepo interfaces.AnalysisOrderRepository, txManager interfaces.TxManager) interfaces.AnalysisOrderUsecase {
	return &AnalysisOrderUsecase{
		analysisRepo:      analysisRepo,
		analysisOrderRepo: analysisOrderRepo,
		txManager:         txManager,
	}
}

func (u *AnalysisOrderUsecase) UpdateAnalysisOrder(ctx context.Context, req *models.UpdateAnalysisOrderRequest) *errors.AppError {
	op := "usecase.AnalysisOrder.UpdateAnalysisOrder"

	ctx, err := u.txManager.Begin(ctx)
	if err != nil {
		return errors.NewDBError(op, err)
	}
	shouldRollback := true
	defer func() {
		if shouldRollback {
			_ = u.txManager.Rollback(ctx)
		}
	}()

	// 1. Проверка заказа
	currentOrder, err := u.analysisOrderRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewNotFoundError(op)
		}
		return errors.NewDBError(op, err)
	}
	if currentOrder.PatientID != req.PatientID {
		return errors.NewForbiddenError(op, "order does not belong to patient")
	}

	// 2. Валидация анализов
	validAnalysisMap, appErr := u.getValidAnalysisMap(ctx)
	if appErr != nil {
		return appErr
	}

	// 3. Подготовка новых item'ов (с логикой!)
	now := time.Now()
	var newItems []entities.AnalysisOrderItem
	for _, dto := range req.OrderItems {
		analysis, exists := validAnalysisMap[dto.AnalysisID]
		if !exists {
			return errors.NewValidationError(op, fmt.Sprintf("analysis %d not found", dto.AnalysisID))
		}

		item := entities.AnalysisOrderItem{
			ID:                dto.ID,
			OrderID:           req.ID,
			AnalysisID:        dto.AnalysisID,
			IsCompleted:       dto.IsCompleted,
			CompletedAt:       dto.CompletedAt,
			PriceAtAssignment: analysis.Price,
			CreatedAt:         now, // ← ЛОГИКА В USECASE!
			UpdatedAt:         now, // ← ЛОГИКА В USECASE!
		}
		// Для существующих item'ов CreatedAt будет перезаписан ниже
		newItems = append(newItems, item)
	}

	// 4. Получение текущих item'ов
	currentItems, err := u.analysisOrderRepo.GetOrderItemsByOrderID(ctx, req.ID)
	if err != nil {
		return errors.NewDBError(op, err)
	}

	// 5. Вычисление изменений
	toCreate, toUpdate, toDelete := u.computeOrderChanges(currentItems, newItems, now)

	// 6. Выполнение операций
	if len(toUpdate) > 0 {
		for _, item := range toUpdate {
			if err := u.analysisOrderRepo.UpdateOrderItem(ctx, item); err != nil {
				return errors.NewDBError(op, err)
			}
		}
	}

	if len(toCreate) > 0 {
		if err := u.analysisOrderRepo.CreateOrderItems(ctx, toCreate); err != nil {
			return errors.NewDBError(op, err)
		}
	}

	if len(toDelete) > 0 {
		if err := u.analysisOrderRepo.DeleteOrderItems(ctx, toDelete); err != nil {
			return errors.NewDBError(op, err)
		}
	}

	shouldRollback = false
	if err := u.txManager.Commit(ctx); err != nil {
		return errors.NewDBError(op, err)
	}

	return nil
}

func (u *AnalysisOrderUsecase) computeOrderChanges(
	currentItems []entities.AnalysisOrderItem,
	newItems []entities.AnalysisOrderItem,
	now time.Time,
) (toCreate, toUpdate []entities.AnalysisOrderItem, toDelete []uint) {
	currentMap := make(map[uint]entities.AnalysisOrderItem)
	for _, item := range currentItems {
		currentMap[item.ID] = item
	}

	for _, newItem := range newItems {
		if newItem.ID == 0 {
			// Новый item
			toCreate = append(toCreate, newItem)
		} else {
			// Существующий item
			if currentItem, exists := currentMap[newItem.ID]; exists {
				newItem.CreatedAt = currentItem.CreatedAt
				newItem.UpdatedAt = now
				toUpdate = append(toUpdate, newItem)
				delete(currentMap, newItem.ID)
			}
		}
	}

	for id := range currentMap {
		toDelete = append(toDelete, id)
	}

	return toCreate, toUpdate, toDelete
}
func (u *AnalysisOrderUsecase) getValidAnalysisMap(ctx context.Context) (map[uint]entities.Analysis, *errors.AppError) {
	op := "usecase.AnalysisOrder.getValidAnalysisMap"

	// Получаем все ID анализов
	analysisIDs, err := u.analysisRepo.GetAllAnalysisIDs(ctx)
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// Получаем полные сущности
	analysisMap := make(map[uint]entities.Analysis)
	for _, id := range analysisIDs {
		analysis, err := u.analysisRepo.GetAnalysisByID(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue // пропускаем несуществующие (маловероятно)
			}
			return nil, errors.NewDBError(op, err)
		}
		analysisMap[id] = *analysis
	}

	return analysisMap, nil
}
