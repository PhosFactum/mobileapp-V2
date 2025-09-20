// internal/domain/entities/invalid_token.go
package entities

import (
	"time"
)

type InvalidToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func (InvalidToken) TableName() string {
	return "invalid_tokens"
}
