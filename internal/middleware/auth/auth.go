package middleware

import (
	"net/http"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		// Проверяем, не в черном списке ли токен
		var invalidToken entities.InvalidToken
		if err := db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&invalidToken).Error; err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalidated"})
			c.Abort()
			return
		}

		// Здесь также должна быть проверка валидности JWT токена
		// (если используете JWT)

		c.Set("token", token)
		c.Next()
	}
}
