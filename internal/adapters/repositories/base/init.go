package base

import (
	"gorm.io/gorm"
)

type contextKey string

const TxContextKey contextKey = "db_transaction"

// BaseRepository — базовый репозиторий с общей логикой
type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}
