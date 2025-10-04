package services

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type contextKey string

const txContextKey contextKey = "db_transaction"

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) interfaces.TxManager {
	return &TxManager{db: db}
}

// Begin начинает транзакцию и кладёт её в контекст
func (tm *TxManager) Begin(ctx context.Context) (context.Context, error) {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return context.WithValue(ctx, txContextKey, tx), nil
}

// Commit коммитит транзакцию из контекста
func (tm *TxManager) Commit(ctx context.Context) error {
	tx := tm.GetTransaction(ctx)
	if tx == nil {
		return nil // нет транзакции — ничего не делаем
	}
	return tx.Commit().Error
}

// Rollback откатывает транзакцию из контекста
func (tm *TxManager) Rollback(ctx context.Context) error {
	tx := tm.GetTransaction(ctx)
	if tx == nil {
		return nil
	}
	return tx.Rollback().Error
}

// getTransaction извлекает транзакцию из контекста
func (tm *TxManager) GetTransaction(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(txContextKey).(*gorm.DB)
	return tx
}
