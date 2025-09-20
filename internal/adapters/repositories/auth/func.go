package auth

import (
	"context"
	"log"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *AuthRepository) GetByLogin(ctx context.Context, login string) (*entities.Doctor, error) {
	op := "repo.Auth.GetByLogin"
	log.Printf("Searching for doctor with login: '%s'", login)

	var doctor entities.Doctor
	if err := r.db.WithContext(ctx).
		Where("phone = ?", login).
		First(&doctor).
		Error; err != nil {

		log.Printf("Error finding doctor: %v", err)
		return nil, errors.NewDBError(op, err)
	}

	log.Printf("Found doctor ID: %d", doctor.ID)
	return &doctor, nil
}

func (r *AuthRepository) InvalidateToken(ctx context.Context, token string) error {
	op := "repo.Auth.InvalidateToken"

	// Создаем запись о недействительном токене
	invalidToken := entities.InvalidToken{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Токен недействителен 24 часа
		CreatedAt: time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(&invalidToken).Error; err != nil {
		log.Printf("Error invalidating token: %v", err)
		return errors.NewDBError(op, err)
	}

	log.Printf("Token invalidated successfully")
	return nil
}
