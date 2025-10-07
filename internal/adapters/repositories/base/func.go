package base

import (
	"context"

	"gorm.io/gorm"
)

// GetDB возвращает транзакцию из контекста или основное подключение
func (br *BaseRepository) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TxContextKey).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return br.db
}
