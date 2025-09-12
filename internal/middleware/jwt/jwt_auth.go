package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secretKey != "your_strong_jwt_secret_here" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization secret"})
			c.Abort()
			return
		}
		// 1. Берём токен из заголовка
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 2. Проверяем Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// 3. Валидируем токен
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что алгоритм совпадает
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 4. Достаём claims и сохраняем в контекст
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		}

		// Пускаем дальше
		c.Next()
	}
}
